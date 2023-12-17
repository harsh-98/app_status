{
  mainnet: {
    aggregatex: [
      // 'https://mainnet.gearbox.foundation/aggregatex/metrics',
      'https://aggregatex.fly.dev/metrics',
    ],
    'third-eye': [
      'https://mainnet.gearbox.foundation/metrics',
    ],
    definder: [
      'https://definder.fly.dev/metrics',
    ],
    'go-liquidator': [
      'https://go-liquidator.fly.dev/metrics',
    ],
    charts_server: [
      'https://mainnet.gearbox.foundation/health',
    ],
    'gearbox-ws-trading': [
      'https://gearbox-ws.fly.dev/metrics',
    ],
    trading_price: [
      // 'https://mainnet.gearbox.foundation/api/tradingview/config',
      'https://trading-price.fly.dev/api/tradingview/config',
    ],
  },
  anvil: {
    'third-eye': [
      'https://testnet.gearbox.foundation/metrics',
    ],
    webhook: [
      'https://testnet.gearbox.foundation/webhook/health',
    ],
    // definder: [
    //   'https://definder.fly.dev/metrics',
    // ],
    // 'go-liquidator': [
    //   'https://go-liquidator.fly.dev/metrics',
    // ],
    charts_server: [
      'https://testnet.gearbox.foundation/health',
    ],
    'gearbox-ws-trading': [
      'https://testnet.gearbox.foundation/gearbox-ws/metrics',
    ],
    gpointbot: [
      'https://testnet.gearbox.foundation/gpointbot/metrics',
    ],
    trading_price: [
      'https://testnet.gearbox.foundation/api/tradingview/config',
    ],
  },
}
