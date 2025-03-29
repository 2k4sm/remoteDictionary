package models

type PutRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type PutResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type GetResponse struct {
	Status  string `json:"status"`
	Key     string `json:"key,omitzero"`
	Value   string `json:"value,omitzero"`
	Message string `json:"message,omitzero"`
}
