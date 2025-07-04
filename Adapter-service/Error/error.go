package Error

import(
	"github.com/lokesh2201013/genet-microservice/Adapter-service/models"
)

func ReturnError(ServiceName string,Err error,Message string)(models.Error){

	return models.Error{
		ServiceName: ServiceName,
		Message:  Message,
		Description: Err.Error(),}
}