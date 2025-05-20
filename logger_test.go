package awslogr_test

import (
	"bytes"
	"context"
	"fmt"
	"log"

	"github.com/aws/smithy-go/logging"
	"github.com/bodgit/awslogr"
	"github.com/tonglil/buflogr"
)

func ExampleLogger() {
	buf := new(bytes.Buffer)

	l, err := awslogr.New(buflogr.NewWithBuffer(buf), awslogr.WithClassificationLevel(func(c logging.Classification) int {
		if c == logging.Debug {
			return 1
		}

		return 0
	}))
	if err != nil {
		log.Fatal(err)
	}

	l.Logf(logging.Warn, "%s", "a warning")

	l, err = awslogr.New(buflogr.NewWithBuffer(buf), awslogr.WithContextKey("ctx"))
	if err != nil {
		log.Fatal(err)
	}

	//nolint:forcetypeassert
	l = l.(logging.ContextLogger).WithContext(context.Background())

	l.Logf(logging.Warn, "%s", "another warning")
	l.Logf(logging.Debug, "%s", "some debug")

	fmt.Print(buf.String())

	// Output:
	// INFO a warning
	// INFO another warning ctx context.Background
	// V[1] some debug ctx context.Background
}
