package operation

import (
	_ "fmt"
	"github.com/gorilla/mux"
	_ "github.com/gorilla/mux"
	_ "log"
	"net/http"
)

func Handler(r *mux.Router) {

	user := r.PathPrefix("/user").Subrouter()
	{
		user.Path("/").Methods(http.MethodPost).HandlerFunc(CreateUser)
		user.Path("/").Methods(http.MethodGet).HandlerFunc(GetAllUser)
		userAction := user.PathPrefix("/{userId}").Subrouter()
		userAction.Path("/").Methods(http.MethodPut).HandlerFunc(UpdateUser)
		userAction.Path("/").Methods(http.MethodDelete).HandlerFunc(UserDelete)
		userAction.Path("/").Methods(http.MethodGet).HandlerFunc(GetUserById)

		address := userAction.PathPrefix("/address").Subrouter()
		{
			address.Path("/").HandlerFunc(CreateAddress).Methods(http.MethodPost)
			address.Path("/").HandlerFunc(DeleteAddressByUserId).Methods(http.MethodDelete)
			address.Path("/user").HandlerFunc(GetAddressByUserId).Methods(http.MethodGet)
			address.Path("/{addressId}").HandlerFunc(GetAddressByAddressId).Methods(http.MethodGet)
			//address.Path("/list").HandlerFunc(GetAllAddress).Methods(http.MethodGet)
			address.Path("/{addressId}").HandlerFunc(UpdateAddress).Methods(http.MethodPut)
			address.Path("/{addressId}").HandlerFunc(DeleteAddressByAddressId).Methods(http.MethodDelete)
		}
		billing := userAction.PathPrefix("/billing").Subrouter()
		{
			billing.Path("/").HandlerFunc(CreateBilling).Methods(http.MethodPost)
			billing.Path("/list").HandlerFunc(GetAllBilling).Methods(http.MethodGet)
			billing.Path("/{billingId}").HandlerFunc(GetBillingById).Methods(http.MethodGet)
			billing.Path("/orderId").HandlerFunc(GetBillingByOrderId).Methods(http.MethodGet)
			billing.Path("/").HandlerFunc(GetBillingByUserId).Methods(http.MethodGet)
			billing.Path("/{billingId}").HandlerFunc(UpdateBilling).Methods(http.MethodPut)
			billing.Path("/{billingId}").HandlerFunc(DeleteBillingById).Methods(http.MethodDelete)
			billing.Path("/").HandlerFunc(DeleteBilling).Methods(http.MethodDelete)
		}
		order := userAction.PathPrefix("/order").Subrouter()
		{
			order.Path("/").HandlerFunc(CreateOrder).Methods(http.MethodPost)
			order.Path("/").HandlerFunc(GetAllOrder).Methods(http.MethodGet)
			order.Path("/user").HandlerFunc(GetOrderByUserId).Methods(http.MethodGet)
			order.Path("/{orderId}").HandlerFunc(GetOrderByOrderId).Methods(http.MethodGet)
			order.Path("/{orderId}").HandlerFunc(UpdateOrder).Methods(http.MethodPut)
			order.Path("/{orderId}").HandlerFunc(DeleteOrderByUserId).Methods(http.MethodDelete)
		}

		booking := userAction.PathPrefix("/booking").Subrouter()
		{
			booking.Path("/").HandlerFunc(CreateBooking).Methods(http.MethodPost)
			booking.Path("/{bookingId}").HandlerFunc(GetBookingByBookingId).Methods(http.MethodGet)
			booking.Path("/").HandlerFunc(GetBookingByUserId).Methods(http.MethodGet)
			booking.Path("/{bookingId}").HandlerFunc(UpdateBooking).Methods(http.MethodPut)
			//booking.Path("/{bookingId").HandlerFunc(DeleteBooking).Methods(http.MethodDelete)
		}
	}
	booking := r.PathPrefix("/booking").Subrouter()
	booking.Path("/").HandlerFunc(GetAllBooking).Methods(http.MethodGet)
	booking.Path("/{bookingId}").HandlerFunc(DeleteBooking).Methods(http.MethodDelete)

	food := r.PathPrefix("/food").Subrouter()
	{
		food.Path("/").HandlerFunc(CreateFood).Methods(http.MethodPost)
		food.Path("/").HandlerFunc(GetAllFood).Methods(http.MethodGet)
		food.Path("/{orderItemId}").HandlerFunc(GetFoodByOrderItemId).Methods(http.MethodGet)
		food.Path("/{foodId}").HandlerFunc(GetFoodById).Methods(http.MethodGet)
		food.Path("/{id}").HandlerFunc(FoodDelete).Methods(http.MethodDelete)
		food.Path("/{id}").HandlerFunc(FoodUpdate).Methods(http.MethodPut)
	}

	table := r.PathPrefix("/table").Subrouter()
	{
		table.Path("/").HandlerFunc(CreateTable).Methods(http.MethodPost)
		table.Path("/").HandlerFunc(GetAllTable).Methods(http.MethodGet)
		table.Path("/{bookingId}").HandlerFunc(GetTableByBookingId).Methods(http.MethodGet)
		table.Path("/{tableId}/table").HandlerFunc(GetTableById).Methods(http.MethodGet)
		table.Path("/{id}").HandlerFunc(UpdateTable).Methods(http.MethodPut)
		table.Path("/{id}").HandlerFunc(DeleteTable).Methods(http.MethodDelete)
	}

	orderItem := r.PathPrefix("/orderItem").Subrouter()
	{
		orderItem.Path("/").HandlerFunc(CreateOrderItem).Methods(http.MethodPost)
		orderItem.Path("/").HandlerFunc(GetAllOrderItem).Methods(http.MethodGet)
		orderItem.Path("/{orderId}").HandlerFunc(GetOrderItemByOrderId).Methods(http.MethodGet)
		orderItem.Path("/{id}").HandlerFunc(GetOrderItemById).Methods(http.MethodGet)
		orderItem.Path("/{id}").HandlerFunc(UpdateOrderItem).Methods(http.MethodPut)
		orderItem.Path("/{id}").HandlerFunc(DeleteOrderItem).Methods(http.MethodDelete)
	}

}
