package db

import (
    "github.com/dimdin/decimal"
    "github.com/jackc/pgx"
    "fmt"
    "errors"
)

//var regexNumber *regexp.Regexp = regexp.MustCompile(`0+\.[0-9]*[1-9][0-9]*$`)

type Number struct {
    decimal.Dec
    Null bool
}

func (p *Number) Scan(vr *pgx.ValueReader) error {
    if vr.Type().DataTypeName != "numeric" {
        return pgx.SerializationError(fmt.Sprintf("Number.Scan cannot decode %s (OID %d)",
        vr.Type().DataTypeName, vr.Type().DataType))
    }

    if vr.Len() == -1 {
        p.SetInt64(0)
        p.Null = true
        return nil
    }

    switch vr.Type().FormatCode {
        case pgx.TextFormatCode:
        s := vr.ReadString(vr.Len())
        err := p.SetString(s)
        if err != nil {
            return pgx.SerializationError(fmt.Sprintf("Received invalid number format: %v", s))
        }
        p.Null = false
        case pgx.BinaryFormatCode:
        return errors.New("binary format not implemented")
        default:
        return fmt.Errorf("unknown format %v", vr.Type().FormatCode)
    }

    return vr.Err()
}