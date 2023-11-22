# Relay 的抽象

抽象分成三个组成

- Request Transformer
  - GRPC Request -> Relay Request
  - OpenAI Request -> Relay Request
- Relay Requester
- Response Transformer
  - Relay Response -> GRPC Response
  - Relay Response -> OpenAI Response