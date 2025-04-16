from pydantic import BaseModel
from typing import List, Optional


class Location(BaseModel):
    city: str
    country: str


class NetWorth(BaseModel):
    estimatedValue: Optional[float]
    currency: str
    source: str


class CareerEvent(BaseModel):
    year: str
    event: str


class Social(BaseModel):
    platform: str
    link: str


class Profile(BaseModel):
    names: List[str]
    gender: str
    dateOfBirth: str
    description: str
    nationality: str
    currentResidence: Location
    netWorth: NetWorth
    industries: List[str]
    occupations: List[str]
    pastOccupations: List[str]
    careerTimeline: List[CareerEvent]
    socials: List[Social]


class Link(BaseModel):
    label: str
    url: str


class Subsidiary(BaseModel):
    name: str
    ownershipPercentage: Optional[float]
    industry: str
    links: List[Link]


class OwnedCompany(BaseModel):
    name: str
    ownershipType: str
    ownershipPercentage: Optional[float]
    industry: str
    status: str
    subsidiaries: List[Subsidiary]
    links: List[Link]


class InvestmentValue(BaseModel):
    value: Optional[float]
    currency: str


class Investment(BaseModel):
    name: str
    type: str
    value: InvestmentValue
    industry: str
    status: str


class FamilyMember(BaseModel):
    name: str
    relationship: str


class Associate(BaseModel):
    name: str
    relationship: str
    associatedCompanies: List[str]


class Source(BaseModel):
    source: str
    confidence: Optional[float]


class ClientProfile(BaseModel):
    profile: Profile
    ownedCompanies: List[OwnedCompany]
    investments: List[Investment]
    family: List[FamilyMember]
    associates: List[Associate]
    sources: List[Source]
