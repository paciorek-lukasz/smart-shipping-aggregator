# Smart Shipping Aggregator

gRPC server aggregating shipping quotes from multiple carriers simultaneously.

## Functionality

Server listens on port `50051` and provides `GetQuotes` method which:
1. Accepts quote request (sender address, recipient address, package)
2. Queries all configured providers in parallel
3. Returns list of offers from various carriers

## Supported Carriers

- DHL
- DPD
- Fedex
- GLS
- Inpost
- UPS

Each provider can be enabled/disabled via environment variables.

## Requirements

- Go 1.25+
- Make

## Installation

```bash
make install-deps
make generate
```

## Configuration

Environment variables (`.env` file in `cmd/server/`):

|Variable|Description|Default|
|--------|------------|-------|
|ENABLE_DHL|Enable DHL provider|false|
|DHL_BASE_URL|DHL API URL|localhost|
|DHL_API_KEY|DHL API key|-|
|ENABLE_DPD|Enable DPD provider|false|
|DPD_BASE_URL|DPD API URL|localhost|
|DPD_API_KEY|DPD API key|-|
|ENABLE_FEDEX|Enable Fedex provider|false|
|FEDEX_BASE_URL|Fedex API URL|localhost|
|FEDEX_API_KEY|Fedex API key|-|
|ENABLE_GLS|Enable GLS provider|false|
|GLS_BASE_URL|GLS API URL|localhost|
|GLS_API_KEY|GLS API key|-|
|ENABLE_INPOST|Enable Inpost provider|false|
|INPOST_BASE_URL|Inpost API URL|localhost|
|INPOST_API_KEY|Inpost API key|-|
|ENABLE_UPS|Enable UPS provider|false|
|UPS_BASE_URL|UPS API URL|localhost|
|UPS_API_KEY|UPS API key|-|
|TIMEOUT|Aggregator timeout in seconds|30s|

## Running

```bash
go run cmd/server/main.go
```

Server starts on `localhost:50051`.

## API

### GetQuotes

gRPC method for fetching shipping quotes.

**Request:**
```protobuf
message GetQuotesRequest {
  Party sender = 1;
  Party recipient = 2;
  Package package = 3;
  DeliveryType delivery_type = 4;
  repeated LocationType location_types = 5;
}

message Party {
  string name = 1;
  Address address = 2;
  string phone = 3;
  string email = 4;
}

message Address {
  string address = 1;
  string postal_code = 2;
  string city = 3;
  string country = 4;
  string longitude = 5;
  string latitude = 6;
}

message Package {
  repeated Item items = 1;
  int32 total_price = 2;
  string currency = 3;
  Dimensions dimensions = 4;
}

message Dimensions {
  int32 length = 1;
  int32 width = 2;
  int32 height = 3;
  float weight = 4;
}

enum DeliveryType {
  DELIVERY_TYPE_UNKNOWN = 0;
  DELIVERY_TYPE_HOME_DELIVERY = 1;
  DELIVERY_TYPE_PICKUP = 2;
}

enum LocationType {
  LOCATION_TYPE_UNKNOWN = 0;
  LOCATION_TYPE_ADDRESS = 1;
  LOCATION_TYPE_PARCEL_LOCKER = 2;
}
```

**Response:**
```protobuf
message GetOptionsResponse {
  repeated Option options = 1;
}

message Option {
  int32 option_id = 1;
  string carrier_product = 2;
  int32 price = 3;
  string currency = 4;
  repeated TimeSlot delivery_time_slots = 5;
  DeliveryType delivery_type = 6;
}
```

### Makefile

```bash
make install-deps  # Install dependencies
make generate      # Generate proto code
make build        # Build project
make run          # Run server
make test        # Run tests
```

### gRPC Reflection

Server supports gRPC Reflection - can be used with grpcurl:

```bash
grpcurl -plaintext localhost:50051 list
grpcurl -plaintext -d '{"sender": {...}}' localhost:50051 smartshippingaggregator.ShippingService/GetQuotes
```

Or in Postman (with gRPC Reflection enabled).

## Architecture

```
┌───────────────┐
│  gRPC Client  │
└──────┬───���────┘
       │
       ▼
┌──────────────┐
│  RPC Handler │
└──────┬────────┘
       │
       ▼
┌──────────────┐
│  Aggregator   │ ──► fan-out to providers
└──────┬────────┘
       │
   ┌───┴───┐
   ▼       ▼
┌──────┐ ┌──────┐
│Provider│ │Provider│ (DHL, DPD, Fedex, GLS, Inpost, UPS)
└──────┘ └──────┘
```

## Resilience

Each request to provider is wrapped with:
- **Circuit Breaker** - protects against cascading failures
- **Timeout** - prevents infinitely blocking

## Project Structure

```
smart-shipping-aggregator/
├── api/shipping/           # Generated gRPC code
├── cmd/server/            # Entry point
│   ├── main.go
│   └── .env              # Configuration
└── internal/
    ├── aggregator/        # Aggregation logic
    ├── config/           # Configuration
    ├── domain/          # Data models
    ├── provider/         # Provider implementations
    │   ├── dhl/
    │   ├── dpd/
    │   ├── fedex/
    │   ├── gls/
    │   ├── inpost/
    │   └── ups/
    ├── resilience/       # Circuit breaker
    └── transport/
        └── rpc/          # gRPC handler
```