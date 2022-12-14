package delivery

import (
	"api-airbnb-alta/features/propertyImage"
	"api-airbnb-alta/middlewares"
	"api-airbnb-alta/utils/helper"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

type PropertyImageDelivery struct {
	propertyImageService propertyImage.ServiceInterface
}

func New(service propertyImage.ServiceInterface, e *echo.Echo) {
	handler := &PropertyImageDelivery{
		propertyImageService: service,
	}

	e.GET("/property_images", handler.GetAll, middlewares.JWTMiddleware())
	e.GET("/property_images/:id", handler.GetById, middlewares.JWTMiddleware())
	e.POST("/property_images", handler.Create, middlewares.JWTMiddleware())
	e.PUT("/property_images/:id", handler.Update, middlewares.JWTMiddleware())
	e.DELETE("/property_images/:id", handler.Delete, middlewares.JWTMiddleware())

	//middlewares.IsAdmin = untuk membatasi akses endpoint hanya admin
	//middlewares.UserOnlySameId = untuk membatasi akses user mengelola data diri sendiri saja

}

func (delivery *PropertyImageDelivery) GetAll(c echo.Context) error {
	query := c.QueryParam("title")
	helper.LogDebug("isi query = ", query)
	results, err := delivery.propertyImageService.GetAll(query)
	if err != nil {
		if strings.Contains(err.Error(), "Get data success. No data.") {
			return c.JSON(http.StatusOK, helper.SuccessWithDataResponse(err.Error(), results))
		}
		return c.JSON(http.StatusBadRequest, helper.FailedResponse(err.Error()))
	}

	dataResponse := fromCoreList(results)

	return c.JSON(http.StatusOK, helper.SuccessWithDataResponse("Success read all data.", dataResponse))
}

func (delivery *PropertyImageDelivery) GetById(c echo.Context) error {
	idParam := c.Param("id")
	id, errConv := strconv.Atoi(idParam)
	if errConv != nil {
		return c.JSON(http.StatusBadRequest, helper.FailedResponse("Error. Id must integer."))
	}
	results, err := delivery.propertyImageService.GetById(id)
	if err != nil {
		if strings.Contains(err.Error(), "Get data success. No data.") {
			return c.JSON(http.StatusOK, helper.SuccessWithDataResponse(err.Error(), results))
		}
		return c.JSON(http.StatusBadRequest, helper.FailedResponse(err.Error()))
	}

	dataResponse := fromCore(results)

	return c.JSON(http.StatusOK, helper.SuccessWithDataResponse("Success read user.", dataResponse))
}

func (delivery *PropertyImageDelivery) Create(c echo.Context) error {
	userInput := InsertRequest{}
	errBind := c.Bind(&userInput) // menangkap data yg dikirim dari req body dan disimpan ke variabel
	if errBind != nil {
		return c.JSON(http.StatusBadRequest, helper.FailedResponse("Error binding data. "+errBind.Error()))
	}

	dataCore := toCore(userInput)
	err := delivery.propertyImageService.Create(dataCore, c)
	if err != nil {
		if strings.Contains(err.Error(), "Error:Field validation") {
			return c.JSON(http.StatusBadRequest, helper.FailedResponse("Some field cannot Empty. Details : "+err.Error()))
		}
		return c.JSON(http.StatusInternalServerError, helper.FailedResponse("Failed insert data. "+err.Error()))
	}
	return c.JSON(http.StatusCreated, helper.SuccessResponse("Success create data"))
}

func (delivery *PropertyImageDelivery) Update(c echo.Context) error {
	idParam := c.Param("id")
	id, errConv := strconv.Atoi(idParam)
	if errConv != nil {
		return c.JSON(http.StatusBadRequest, helper.FailedResponse("Error. Id must integer."))
	}

	userInput := UpdateRequest{}
	errBind := c.Bind(&userInput) // menangkap data yg dikirim dari req body dan disimpan ke variabel
	if errBind != nil {
		return c.JSON(http.StatusBadRequest, helper.FailedResponse("Error binding data. "+errBind.Error()))
	}

	// validasi data di proses oleh user ybs
	userId := middlewares.ExtractTokenUserId(c)
	if userId < 1 {
		return c.JSON(http.StatusBadRequest, helper.FailedResponse("Failed load user id from JWT token, please check again."))
	}
	propertyImageData, errGet := delivery.propertyImageService.GetById(id)
	if errGet != nil {
		return c.JSON(http.StatusBadRequest, helper.FailedResponse(errGet.Error()))
	}

	propertyData, errGet := delivery.propertyImageService.GetPropertyById(int(propertyImageData.PropertyID))
	if errGet != nil {
		return c.JSON(http.StatusBadRequest, helper.FailedResponse(errGet.Error()))
	}

	fmt.Println("\n\nProp image data = ", propertyImageData)
	fmt.Println("\n\nid = ", userId, " = prop user id =", propertyData.UserID)

	if userId != int(propertyData.UserID) {
		return c.JSON(http.StatusBadRequest, helper.FailedResponse("Failed process data, data must be yours."))
	}

	dataCore := toCore(userInput)
	err := delivery.propertyImageService.Update(dataCore, id, c)
	if err != nil {
		if strings.Contains(err.Error(), "Error:Field validation") {
			return c.JSON(http.StatusBadRequest, helper.FailedResponse("Some field cannot Empty. Details : "+err.Error()))
		}
		return c.JSON(http.StatusInternalServerError, helper.FailedResponse("Failed update data. "+err.Error()))
	}

	return c.JSON(http.StatusCreated, helper.SuccessResponse("Success update data."))
}

func (delivery *PropertyImageDelivery) Delete(c echo.Context) error {
	idParam := c.Param("id")
	id, errConv := strconv.Atoi(idParam)
	if errConv != nil {
		return c.JSON(http.StatusBadRequest, helper.FailedResponse("Error. Id must integer."))
	}

	// validasi data di proses oleh user ybs
	userId := middlewares.ExtractTokenUserId(c)
	if userId < 1 {
		return c.JSON(http.StatusBadRequest, helper.FailedResponse("Failed load user id from JWT token, please check again."))
	}
	propertyImageData, errGet := delivery.propertyImageService.GetById(id)
	if errGet != nil {
		return c.JSON(http.StatusBadRequest, helper.FailedResponse(errGet.Error()))
	}

	propertyData, errGet := delivery.propertyImageService.GetPropertyById(int(propertyImageData.PropertyID))
	if errGet != nil {
		return c.JSON(http.StatusBadRequest, helper.FailedResponse(errGet.Error()))
	}

	fmt.Println("\n\nProp image data = ", propertyImageData)
	fmt.Println("\n\nid = ", userId, " = prop user id =", propertyData.UserID)

	if userId != int(propertyData.UserID) {
		return c.JSON(http.StatusBadRequest, helper.FailedResponse("Failed process data, data must be yours."))
	}

	err := delivery.propertyImageService.Delete(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.FailedResponse(err.Error()))
	}

	return c.JSON(http.StatusOK, helper.SuccessResponse("Success delete data."))
}
