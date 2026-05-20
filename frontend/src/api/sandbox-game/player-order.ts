import { request } from '@/api/http'
import type {
  DeliverOrdersRequest,
  DeliverOrdersResult,
  PassOrderSegmentRequest,
  PassOrderSegmentResult,
  PlayerOrderYearView,
  SelectOrderRequest,
  SelectOrderResult,
  SubmitMarketInvestmentRequest,
  SubmitMarketInvestmentResult,
} from '@/types/sandbox-game-order'

export function getPlayerOrderYearView(yearNo: number) {
  return request<PlayerOrderYearView>(`/api/v1/sandbox-game/player-order/get-year-view?yearNo=${yearNo}`)
}

export function submitPlayerMarketInvestment(payload: SubmitMarketInvestmentRequest) {
  return request<SubmitMarketInvestmentResult>('/api/v1/sandbox-game/player-order/submit-market-investments', {
    method: 'POST',
    body: JSON.stringify(payload),
  })
}

export function selectPlayerOrder(payload: SelectOrderRequest) {
  return request<SelectOrderResult>('/api/v1/sandbox-game/player-order/select-order', {
    method: 'POST',
    body: JSON.stringify(payload),
  })
}

export function passPlayerOrderSegment(payload: PassOrderSegmentRequest) {
  return request<PassOrderSegmentResult>('/api/v1/sandbox-game/player-order/pass-segment', {
    method: 'POST',
    body: JSON.stringify(payload),
  })
}

export function deliverPlayerOrders(payload: DeliverOrdersRequest) {
  return request<DeliverOrdersResult>('/api/v1/sandbox-game/player-order/deliver-orders', {
    method: 'POST',
    body: JSON.stringify(payload),
  })
}
