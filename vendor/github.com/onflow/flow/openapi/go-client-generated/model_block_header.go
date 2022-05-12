/*
 * Access API
 *
 * No description provided (generated by Swagger Codegen https://github.com/swagger-api/swagger-codegen)
 *
 * API version: 1.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package swagger
import (
	"time"
)

type BlockHeader struct {
	Id string `json:"id"`
	ParentId string `json:"parent_id"`
	Height string `json:"height"`
	Timestamp time.Time `json:"timestamp"`
	ParentVoterSignature string `json:"parent_voter_signature"`
}
