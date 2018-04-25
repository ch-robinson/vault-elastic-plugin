package elastic

import (
	"testing"

	"github.com/ch-robinson/vault-elastic-plugin/testdata"
	"github.com/hashicorp/vault/builtin/logical/database/dbplugin"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	mockHTTP := testdata.NewMockHTTP(nil)
	httpClient := testdata.NewMockHTTPClient(nil, mockHTTP)

	dbType, err := New(httpClient)

	assert.Nil(t, err)

	dbMiddleware := dbType.(*dbplugin.DatabaseErrorSanitizerMiddleware)

	name, err := dbMiddleware.Type()

	assert.Nil(t, err)
	assert.Equal(t, name, "elasticdb")
}
