package v1

import "fmt"

type RabbitConfig struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Host      string `json:"host"`
	Port      string `json:"port"`
	QueueName string `json:"queue_name"`
}

func (r *RabbitConfig) Validate() error {

	b := (r.Username == "") &&
		(r.Password == "") &&
		(r.Host == "") &&
		(r.Port == "") &&
		(r.QueueName == "")

	if b {
		return fmt.Errorf("[RABBIT] some configs are empty")
	}

	return nil
}
