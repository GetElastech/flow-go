/*
 * Access API
 *
 * No description provided (generated by Swagger Codegen https://github.com/swagger-api/swagger-codegen)
 *
 * API version: 1.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package generated

type Event struct {
	Type_ string `json:"type"`

	TransactionId string `json:"transaction_id,omitempty"`

	TransactionIndex string `json:"transaction_index,omitempty"`

	EventIndex string `json:"event_index,omitempty"`

	Payload string `json:"payload,omitempty"`
}
