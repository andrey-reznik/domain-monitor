package service

import (
	"errors"
	"golang.org/x/net/idna"
	"log"

	"github.com/nwesterhausen/domain-monitor/configuration"
)

type ServicesWhois struct {
	store configuration.WhoisCacheStorage
}

func NewWhoisService(store configuration.WhoisCacheStorage) *ServicesWhois {
	return &ServicesWhois{store: store}
}

func (s *ServicesWhois) GetWhois(fqdn string) (configuration.WhoisCache, error) {
	punycodeFqdn, _ := idna.ToASCII(fqdn)
	for _, entry := range s.store.FileContents.Entries {
		if entry.FQDN == punycodeFqdn {
			return entry, nil
		}
	}
	log.Println("🙅 WHOIS entry cache miss for", punycodeFqdn)

	// Since we cache missed, let's try to fetch the WHOIS entry instead
	s.store.Add(punycodeFqdn)
	// Try to get the entry again
	for _, entry := range s.store.FileContents.Entries {
		if entry.FQDN == punycodeFqdn {
			return entry, nil
		}
	}

	return configuration.WhoisCache{}, errors.New("entry missing")
}

func (s *ServicesWhois) MarkAlertSent(fqdn string, alert configuration.Alert) bool {
	punycodeFqdn, _ := idna.ToASCII(fqdn)
	for i := range s.store.FileContents.Entries {
		if s.store.FileContents.Entries[i].FQDN == punycodeFqdn {
			s.store.FileContents.Entries[i].MarkAlertSent(alert)
			return true
		}
	}
	return false
}
