package jwt

import (
	"errors"
	"fmt"
	"strings"

	"github.com/nats-io/jwt/v2"
)

var (
	ErrNotSelfSigned  = errors.New("jwt is not self-signed")
	ErrInvalidAccount = errors.New("jwt is invalid")
)

func Validate(claimJWT string) error {
	appClaims, err := jwt.DecodeAccountClaims(claimJWT)
	if err != nil {
		return err
	}
	if !appClaims.IsSelfSigned() {
		return ErrNotSelfSigned
	}
	vr := jwt.CreateValidationResults()
	appClaims.Validate(vr)

	errs := blockingIssues(vr)
	if len(errs) > 0 {
		return invalidAccountError(errs)
	}

	if err := ValidateAppClaims(appClaims); err != nil {
		return err
	}

	return nil
}

// ValidateAppClaims checks imports and exports of self-signed claims to have necessary permissions.
func ValidateAppClaims(accClaims *jwt.AccountClaims) error {
	if accClaims == nil || !accClaims.IsSelfSigned() {
		return nil
	}

	if len(accClaims.Imports) == 0 && len(accClaims.Exports) == 0 {
		return ErrNoImportAndExport
	}

	isValidSubjectName := func(subj string) bool {
		s := string(subj)
		return !strings.HasPrefix(s, ">") && !strings.HasPrefix(s, "*")
	}

	for _, imp := range accClaims.Imports {
		if !isValidSubjectName(string(imp.Subject)) {
			return fmt.Errorf("permissions violation for Import of %q", imp.Subject)
		}
		if !accClaims.DefaultPermissions.Sub.Allow.Contains(string(imp.Subject)) {
			return fmt.Errorf("permissions violation for Import of %q", imp.Subject)
		}
	}
	for _, exp := range accClaims.Exports {
		if !isValidSubjectName(string(exp.Subject)) {
			return fmt.Errorf("permissions violation for Export of %q", exp.Subject)
		}
		if !accClaims.DefaultPermissions.Pub.Allow.Contains(string(exp.Subject)) {
			return fmt.Errorf("permissions violation for Export of %q", exp.Subject)
		}
	}
	return nil
}

func blockingIssues(vr *jwt.ValidationResults) (issues []error) {
	for _, i := range vr.Issues {
		if i.Blocking {
			issues = append(issues, i)
		}

		if i.TimeCheck {
			issues = append(issues, i)
		}
	}
	return issues
}

func invalidAccountError(errs []error) error {
	if len(errs) == 0 {
		return nil
	}

	var sb strings.Builder
	for i := 0; i < len(errs)-1; i++ {
		sb.WriteString(errs[i].Error())
		sb.WriteString("; ")
	}
	sb.WriteString(errs[len(errs)-1].Error())
	return fmt.Errorf("%w: %s", ErrInvalidAccount, sb.String())
}
