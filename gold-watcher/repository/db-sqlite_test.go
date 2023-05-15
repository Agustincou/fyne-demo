package repository

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

//WARNING: coupled tests. run in order to succeed.
//ToDo: improve tests to uncouple

func Test_Migrate(t *testing.T) {
	err := testSqliteRepo.Migrate()

	assert.Nil(t, err)
}

func Test_Insert(t *testing.T) {
	h := Holding{
		Amount:        2,
		PurchaseDate:  time.Now(),
		PurchasePrice: 125,
	}

	result, err := testSqliteRepo.InsertHolding(h)

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.NotEqual(t, int64(0), result.ID)
}

func Test_GetAllHoldings(t *testing.T) {
	h, err := testSqliteRepo.AllHoldings()

	assert.Nil(t, err)
	assert.Equal(t, 1, len(h))
}

func Test_GetHoldingByID(t *testing.T) {
	h, err := testSqliteRepo.GetHoldingByID(1)

	assert.Nil(t, err)
	assert.Equal(t, int(2), h.Amount)
}

func Test_GetHoldingByNotExistingID(t *testing.T) {
	h, err := testSqliteRepo.GetHoldingByID(1234567)

	assert.NotNil(t, err)
	assert.Equal(t, "sql: no rows in result set", err.Error())
	assert.Nil(t, h)
}

func Test_Update(t *testing.T) {
	h, _ := testSqliteRepo.GetHoldingByID(1)

	assert.Equal(t, 125, h.PurchasePrice) //Created by Test_Insert()
	assert.Equal(t, 2, h.Amount)          //Created by Test_Insert()

	err := testSqliteRepo.UpdateHolding(1, Holding{
		Amount:        5,
		PurchaseDate:  time.Now(),
		PurchasePrice: 999,
	})

	h, _ = testSqliteRepo.GetHoldingByID(1)

	assert.Nil(t, err)
	assert.Equal(t, 999, h.PurchasePrice)
	assert.Equal(t, 5, h.Amount)
}

func Test_Delete(t *testing.T) {
	err := testSqliteRepo.DeleteHolding(1)

	assert.Nil(t, err)

	_, err = testSqliteRepo.GetHoldingByID(1)

	assert.NotNil(t, err)
	assert.Equal(t, "sql: no rows in result set", err.Error())

	err = testSqliteRepo.DeleteHolding(1)

	assert.NotNil(t, err)
	assert.Equal(t, "delete failed", err.Error())
}
