package alditalk

/*

IN COOKIES STEHT DIESE customer ID
die brauchen für selfcare-dashboard api

/scs/bff/scs-209-selfcare-dashboard-bff/selfcare-dashboard/v1/offers/C-0005150729
{
  "subscribedOffers": [
    {
      "offerName": "Tarif M",
      // ... viele andere Eigenschaften ...
      "pack": [
        {
          "unit": "kilobytes",
          "nextExpirationDate": "2025-09-14T09:46:19Z",
          "tariff": "14.99",
          "used": "19186590",       // <--- DIESER WERT
          "type": "data",
          "allocated": "31457280",
          "balanceAttributeReference": "dataGrantAmount"
        },
        {
          // ... das zweite Datenpaket (für EU-Roaming) ...
          "used": "0"
        }
      ]
    }
  ]
}
*/

// packs abfragen
func GetPacks(customerID string) {}

// packs nachbuchen
func ExtendPack(customerID string) {}
