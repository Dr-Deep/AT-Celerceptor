package atfram

/*

logged in as : A-5138413 : !

/*
 ? Selfcare Dashboard

 * CUSTOMER API: persönliche Daten, MSISDN, Vertrag, ICCID usw
 GET /scs/bff/scs-209-selfcare-dashboard-bff/selfcare-dashboard/v1/customer

 * ACCOUNT OVERVIEW API: Verbrauch (Daten, SMS, Minuten), Tarif, BillingAccountId, Produktdaten
 GET /scs/bff/scs-209-selfcare-dashboard-bff/selfcare-dashboard/v1/account-overview/content?contractId=04913887&subscriberType=Prepaid

 * OFFERS API: aktueller Tarif, gebuchte Optionen, nachbuchbare Addons
 GET /scs/bff/scs-209-selfcare-dashboard-bff/selfcare-dashboard/v1/offers/{tef_customer_id}?warningDays=28&contractId=04913887&productType=Mobile_Product_


/scs/bff/scs-209-selfcare-dashboard-bff/selfcare-dashboard/v1/customer
/scs/bff/scs-209-selfcare-dashboard-bff/selfcare-dashboard/v1/account-overview/content
/scs/bff/scs-209-selfcare-dashboard-bff/selfcare-dashboard/v1/offers/{customerId}


* https://www.alditalk-kundenportal.de/scs/bff/scs-209-selfcare-dashboard-bff/selfcare-dashboard/v1/customer
* https://www.alditalk-kundenportal.de/scs/bff/scs-207-customer-master-data-bff/customer-master-data/v1/navigation-list?msisdn=04913887
* https://www.alditalk-kundenportal.de/scs/bff/scs-209-selfcare-dashboard-bff/selfcare-dashboard/v1/offers/C-0005150729?warningDays=28&contractId=04913887&productType=Mobile_Product_Offer
* https://www.alditalk-kundenportal.de/content/alditalk/de/de/user/auth/account-overview.webtracking.json?timeInterval=&tef_customer_id=C-0005150729
* https://www.alditalk-kundenportal.de/scs/bff/scs-207-customer-master-data-bff/customer-master-data/v1/navigation-list?msisdn=04913887
* https://www.alditalk-kundenportal.de/scs/bff/scs-209-selfcare-dashboard-bff/selfcare-dashboard/v1/account-overview/content?contractId=04913887&subscriberType=Prepaid

* https://www.alditalk-kundenportal.de/etc.clientlibs/alditalk/components/page/clientlib-data.lc-af94a09d05186552191da7d14b930e0e-lc.min.js

? https://www.alditalk-kundenportal.de/portal/auth/uebersicht/
? https://www.alditalk-kundenportal.de/user/auth/account-overview/

* https://www.alditalk-kundenportal.de/scs/bff/scs-209-selfcare-dashboard-bff/selfcare-dashboard/v1/offers/C-0005150729?warningDays=28&contractId=04913887&productType=Mobile_Product_Offer

/scs/bff/scs-209-selfcare-dashboard-bff/selfcare-dashboard/v1/offers/C-0005150729
*/

/*
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
