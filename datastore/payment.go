package datastore

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/google/uuid"
	"payments/currency"
)

type PaymentStatus string

const (
	PaymentInitiated = "initiated"
	PaymentFailed    = "failed"
	PaymentExpired   = "expired"
	PaymentPaid      = "paid"
	PaymentRefunded  = "refunded"
)

type Payment struct {
	ID         uuid.UUID
	ExternalID string
	Status     PaymentStatus
	Amount     currency.Amount
}

// InMemoryPaymentRepository stores all the payments in the memory.
// In real life we should persist all the payments in the DB.
type InMemoryPaymentRepository struct {
	payments map[uuid.UUID]Payment
	locker   *sync.RWMutex
}

func NewInMemoryPaymentRepository() *InMemoryPaymentRepository {
	return &InMemoryPaymentRepository{
		payments: make(map[uuid.UUID]Payment),
		locker:   &sync.RWMutex{}, // RWMutex is not really required, just for the exercise it's being used to show the possible edge cases
	}
}

func (i *InMemoryPaymentRepository) Create(_ context.Context, p Payment) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("InMemoryPaymentRepository.UpdateInitiatedByExternalID(%+q): %w", p.ID, err)
		}
	}()

	err = func() error {
		i.locker.RLock()
		defer i.locker.RUnlock()

		if _, ok := i.payments[p.ID]; ok {
			return errors.New("payment already exists")
		}

		for _, x := range i.payments {
			if x.ID == p.ID || x.ExternalID == p.ExternalID {
				return errors.New("payment already exists")
			}
		}

		return nil
	}()
	if err != nil {
		return err
	}

	func() {
		i.locker.Lock()
		defer i.locker.Unlock()

		i.payments[p.ID] = p

		// TODO in real life the logger would be injected, and most likely would not be used in the repository.
		// since it's for mocking purposes only, I'm logging the value here
		log.Default().Println(fmt.Sprintf("Created payment %+q for amount %s, external_id=%+q", p.ID, p.Amount, p.ExternalID))
	}()

	return nil
}

func (i *InMemoryPaymentRepository) UpdateInitiatedByExternalID(_ context.Context, extID string, status PaymentStatus) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("InMemoryPaymentRepository.UpdateInitiatedByExternalID(%+q): %w", extID, err)
		}
	}()

	i.locker.Lock()
	defer i.locker.Unlock()

	for _, x := range i.payments {
		if x.ExternalID != extID {
			continue
		}

		if x.Status != PaymentInitiated {
			return fmt.Errorf("payment has status %+q", x.Status)
		}

		x.Status = status

		i.payments[x.ID] = x

		return nil
	}

	return fmt.Errorf("payment does not exist")
}

func (i *InMemoryPaymentRepository) GetByID(_ context.Context, id uuid.UUID) (Payment, error) {
	var err error

	p, ok := i.payments[id]
	if !ok {
		err = fmt.Errorf("payment %+q does not exist", p.ID.String())
	}

	return p, err
}

func (i *InMemoryPaymentRepository) RefundByID(_ context.Context, paymentID uuid.UUID) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("InMemoryPaymentRepository.RefundByID(%+q): %w", paymentID, err)
		}
	}()

	i.locker.Lock()
	defer i.locker.Unlock()

	for _, x := range i.payments {
		if x.ID != paymentID {
			continue
		}

		if x.Status != PaymentPaid {
			return fmt.Errorf("payment has status %+q", x.Status)
		}

		x.Status = PaymentRefunded

		i.payments[x.ID] = x

		return nil
	}

	return fmt.Errorf("payment does not exist")
}
