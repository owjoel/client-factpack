from pydantic import BaseModel, Field
from typing import List, Optional

class Residence(BaseModel):
    city: str
    country: str


class NetWorth(BaseModel):
    estimated_value: int = Field(..., alias="estimatedValue")
    currency: str
    source: str


class SocialMedia(BaseModel):
    platform: str
    username: str


class Contact(BaseModel):
    work_address: str = Field(..., alias="workAddress")
    phone: str


class InvestmentValue(BaseModel):
    value: int
    currency: str


class Investment(BaseModel):
    name: str
    type: str
    value: InvestmentValue
    industry: str
    status: str
    source: str


class Associate(BaseModel):
    name: str
    relationship: str
    associated_companies: List[str] = Field(..., alias="associatedCompanies")


class Metadata(BaseModel):
    sources: List[str]


class Profile(BaseModel):
    name: str
    age: int
    nationality: str
    current_residence: Residence = Field(..., alias="currentResidence")
    net_worth: NetWorth = Field(..., alias="netWorth")
    industries: List[str]
    occupations: List[str]
    socials: List[SocialMedia]
    contact: Contact


class Client(BaseModel):
    profile: Profile
    investments: Optional[List[Investment]]
    associates: Optional[List[Associate]]
    metadata: Metadata