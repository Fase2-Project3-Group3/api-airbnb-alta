package repository

import (
	"api-airbnb-alta/features/booking"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type bookingRepository struct {
	db *gorm.DB
}

func New(db *gorm.DB) booking.RepositoryInterface {
	return &bookingRepository{
		db: db,
	}
}

// Create implements booking.RepositoryInterface
func (repo *bookingRepository) Create(input booking.Core) error {
	var properties Property

	tx1 := repo.db.First(&properties, input.PropertyID)
	if tx1.Error != nil {
		return tx1.Error
	}

	input.PricePerNight = properties.PricePerNight
	// input.GrossAmount = properties.PricePerNight
	input.BookingStatus = "Complete_payment"

	// kalkuasi jumlah har menginap
	stayDay := input.CheckoutDate.Sub(input.CheckinDate)
	// kalkulasi gross_amout
	input.GrossAmount = float64(stayDay.Hours()/24) * properties.PricePerNight

	fmt.Println("\n\nIsi input ", input)

	dataGorm := fromCore(input)
	tx := repo.db.Create(&dataGorm)
	if tx.Error != nil {
		return errors.New("failed create data")
	}

	if tx.RowsAffected == 0 {
		fmt.Println("Failed Create Data Booking. Row Affected ", tx.RowsAffected)
		return errors.New("failed create data")
	}

	fmt.Println("\n\n\ncreate Affected ", tx.RowsAffected)

	fmt.Println("\n\nHasil insert booking ", dataGorm)

	// var id int
	// ty := repo.db.Raw("SELECT LAST_INSERT_ID()").Scan(&id)
	// if ty.Error != nil {
	// 	return errors.New("failed create data")
	// }

	// if ty.RowsAffected == 0 {
	// 	fmt.Println("Failed select last insert. Row Affected ", ty.RowsAffected)
	// 	return errors.New("failed create data")
	// }

	// fmt.Println("Gross Row Affected ", ty.RowsAffected)

	// ta := repo.db.Exec("UPDATE bookings SET gross_amount = (SELECT DATEDIFF(checkout_date, checkin_date) *price_per_night) WHERE id = ?", id)
	// if ta.Error != nil {
	// 	return errors.New("failed create data")
	// }

	// if ta.RowsAffected == 0 {
	// 	fmt.Println("Failed update gross ammount. Row Affected ", ta.RowsAffected)
	// 	return errors.New("failed create data")
	// }

	// fmt.Println("Gross Row Affected ", ta.RowsAffected)

	return nil

	// var properties Property
	// // var newData Booking
	// tx1 := repo.db.First(&properties, input.PropertyID)
	// if tx1.Error != nil {
	// 	return tx1.Error
	// }

	// input.PricePerNight = properties.PricePerNight
	// // input.GrossAmount = properties.PricePerNight
	// input.BookingStatus = "Complete_payment"
	// input.Property.PropertyName = properties.PropertyName

	// userGorm := fromCore(input)
	// tx := repo.db.Create(&userGorm) // proses insert data
	// if tx.Error != nil {
	// 	return tx.Error
	// }
	// if tx.RowsAffected == 0 {
	// 	return errors.New("insert failed")
	// }

	// //  mencari jumlah hari
	// // tx2 := repo.db.Raw("SELECT DATEDIFF(checkout_date, checkin_date) from bookings where id = (select max(id) from bookings)")
	// // if tx2.Error != nil {
	// // 	return tx.Error
	// // }
	// // grossAmount := properties.PricePerNight * tx2
	// // tx3 := repo.db.Model(&newData).Where("id = ?", "last_insert_id()").Update("gross_amount", grossAmount)
	// // if tx3.Error != nil {
	// // 	return tx.Error
	// }
	// return nil
}

// GetAll implements booking.RepositoryInterface
func (repo *bookingRepository) GetAll(userId int) (data []booking.Core, err error) {
	var results []Booking

	tx := repo.db.Preload("User").Preload("Property").Where("user_id = ?", userId).Find(&results)
	if tx.Error != nil {
		return nil, tx.Error
	}
	var dataCore = toCoreList(results)
	return dataCore, nil
}

// GetById implements booking.RepositoryInterface
func (repo *bookingRepository) GetById(id int, userId int) (data booking.Core, err error) {
	var result Booking

	tx := repo.db.Preload("User").Preload("Property").Where("user_id = ?", userId).First(&result, id)

	if tx.Error != nil {
		return data, tx.Error
	}

	if tx.RowsAffected == 0 {
		return data, tx.Error
	}

	var dataCore = result.toCore()
	return dataCore, nil
}

func (repo *bookingRepository) GetAvailability(propertyId uint, checkinDate time.Time, checkoutData time.Time) (result string, err error) {
	var properties []Property
	queryBuilder := fmt.Sprintf("SELECT * FROM bookings WHERE property_id = %d AND (('%s' BETWEEN checkin_date AND checkout_date) OR ('%s' BETWEEN checkin_date AND checkout_date));", propertyId, checkinDate, checkoutData)

	fmt.Println("\n\n query ", queryBuilder)

	tx := repo.db.Raw(queryBuilder).Find(&properties)

	if tx.Error != nil {
		return "Not Available", tx.Error
	}

	affectedRow := tx.RowsAffected
	fmt.Println("\n\nHasil check availbility, \n Checkin date = ", checkinDate, " \n Checkout date = ", checkoutData, " \n Hasil Row = ", affectedRow)

	if affectedRow == 0 {
		return "Available", nil
	}

	return "Not Available", nil
}

// GetById implements property.RepositoryInterface
func (repo *bookingRepository) GetPropertyById(id int) (affectedRow int, err error) {
	var property Property

	tx := repo.db.First(&property, id)

	if tx.Error != nil {
		return 0, tx.Error
	}

	if tx.RowsAffected == 0 {
		return 0, tx.Error
	}

	return int(tx.RowsAffected), nil
}
