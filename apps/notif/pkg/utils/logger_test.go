package utils_test

import (
	"os"
	"testing"

	"github.com/owjoel/client-factpack/apps/notif/pkg/utils"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestInitLogger(t *testing.T) {
	utils.InitLogger()

	// Assert logger is not nil
	assert.NotNil(t, utils.Logger)

	// Check logger output (should be os.Stdout)
	assert.Equal(t, os.Stdout, utils.Logger.Out)

	// Check logger formatter type
	_, ok := utils.Logger.Formatter.(*logrus.JSONFormatter)
	assert.True(t, ok, "Formatter should be JSONFormatter")

	// Check log level
	assert.Equal(t, logrus.InfoLevel, utils.Logger.Level)
}
