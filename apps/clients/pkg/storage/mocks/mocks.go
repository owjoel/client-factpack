package mocks

var MockClientJSON = `{
  "_id": {
    "$oid": "67cad1efba24c294b3389cea"
  },
  "profile": {
    "name": "Filippo Ghirelli",
    "age": 38,
    "description": "Italian entrepreneur and green energy investor",
    "nationality": "Italian",
    "currentResidence": {
      "city": "Milan",
      "country": "Italy"
    },
    "netWorth": {
      "estimatedValue": 15000000,
      "currency": "EUR",
      "source": "Forbes Europe 2023"
    },
    "industries": [
      "Green Energy",
      "Renewable Resources"
    ],
    "occupations": [
      "Entrepreneur",
      "Investor"
    ],
    "socials": [
      {
        "platform": "LinkedIn",
        "username": "filippo-ghirelli"
      },
      {
        "platform": "Twitter",
        "username": "FilippoGhirelli"
      }
    ],
    "contact": {
      "workAddress": "Via della Repubblica, 5, Milan, Italy",
      "phone": "+39 02 12345678"
    }
  },
  "investments": [
    {
      "name": "Solar Innovations Inc.",
      "type": "Equity",
      "value": {
        "value": 5000000,
        "currency": "EUR"
      },
      "industry": "Solar Energy",
      "status": "Active"
    },
    {
      "name": "Greener Technologies Ltd.",
      "type": "Debt",
      "value": {
        "value": 3000000,
        "currency": "EUR"
      },
      "industry": "Renewable Technology",
      "status": "Active"
    }
  ],
  "associates": [
    {
      "name": "Giovanni Rossi",
      "relationship": "Business Partner",
      "associatedCompanies": [
        "Solar Innovations Inc.",
        "Greener Technologies Ltd."
      ]
    },
    {
      "name": "Laura Bianchi",
      "relationship": "Investment Advisor",
      "associatedCompanies": [
        "Bianchi Investments"
      ]
    }
  ]
}`
