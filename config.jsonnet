{
  mainnet: {
    aggregatex: [
      'https://testnet.gearbox.foundation/aggregatex/metrics',
    ],
    'third-eye': [
      'https://testnet.gearbox.foundation/mainnet/metrics',
    ],
    // definder: [
    //   'https://definder.fly.dev/metrics',
    // ],
    'liquidator-v2': [
      'https://go-liquidator.fly.dev/metrics',
    ],
    'liquidator-v3': [
      'https://liquidator-v3.fly.dev/metrics',
    ],
    charts_server: [
      // 'https://mainnet.gearbox.foundation/health',
      'https://charts-server.fly.dev/health',
    ],
    'gearbox-ws': [
      'https://gearbox-ws.fly.dev/metrics',
    ],
    gpointbot: [
      'https://gpointbot.fly.dev/metrics',
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
    'gearbox-ws': [
      'https://testnet.gearbox.foundation/gearbox-ws/metrics',
    ],
    gpointbot: [
      'https://testnet.gearbox.foundation/gpointbot/metrics',
    ],
    trading_price: [
      'https://testnet.gearbox.foundation/api/tradingview/config',
    ],
    'http-logger': [
      // 'https://mainnet.gearbox.foundation/api/tradingview/config',
      'https://testnet.gearbox.foundation/logger/metrics',
    ],
  },
}
