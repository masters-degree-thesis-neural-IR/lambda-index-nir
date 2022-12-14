package memory

import (
	"context"
	"encoding/json"
	"google.golang.org/grpc"
	"lambda-index-nir/service/application/exception"
	"lambda-index-nir/service/application/repositories"
	"lambda-index-nir/service/infraestructure/speedup"
	"log"
	"time"
)

type SpeedupRepository struct {
}

func NewSpeedupRepository() repositories.IndexMemoryRepository {

	return &SpeedupRepository{}
}

func (r *SpeedupRepository) FindByTerm(term string) ([]string, error) {

	log.Println("FindByTerm")

	var conn *grpc.ClientConn
	conn, err := grpc.Dial("172.31.2.165:9000", grpc.WithInsecure())

	if err != nil {
		log.Println("Tem erro aqui")
		log.Println(err.Error())
		return nil, exception.ThrowValidationError("Not is possible connect to RCP Server.")
	}
	defer conn.Close()

	client := speedup.NewDataServiceClient(conn)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	response, err := client.GetData(ctx, &speedup.RequestDataKey{
		Key: term,
	})

	log.Println("Teve resposta?")
	log.Println(response.GetValue())

	if ctx.Err() == context.Canceled {
		return nil, exception.ThrowValidationError("RPC Client cancelled, abandoning.")
	}

	if err != nil {
		log.Println("Tem erro aqui")
		log.Println(err.Error())
		return nil, err
	}

	if response.GetException() != "" {
		log.Println(err.Error())
		return nil, exception.ThrowValidationError(response.GetException())
	}

	var locDocuments []string
	json.Unmarshal([]byte(response.GetValue()), &locDocuments)
	return locDocuments, nil

}

func (r *SpeedupRepository) Update(term string, documents []string) error {
	return r.Save(term, documents)
}

func (r *SpeedupRepository) Save(term string, documents []string) error {

	log.Println("Save")

	var conn *grpc.ClientConn
	conn, err := grpc.Dial("172.31.2.165:9000", grpc.WithInsecure())
	if err != nil {
		log.Println(err.Error())
		return exception.ThrowValidationError("Not is possible connect to RCP Server.")
	}
	defer conn.Close()

	client := speedup.NewDataServiceClient(conn)

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	value, err := json.Marshal(documents)
	response, err := client.SetData(ctx, &speedup.RequestDataKeyValue{
		Key:   term,
		Value: string(value),
	})

	if ctx.Err() == context.Canceled {
		log.Println(err.Error())
		return exception.ThrowValidationError("RPC Client cancelled, abandoning.")
	}

	if err != nil {
		return err
	}

	if response.GetException() != "" {
		log.Println(err.Error())
		return exception.ThrowValidationError(response.GetException())
	}

	return nil

}
