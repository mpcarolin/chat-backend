import '@testing-library/jest-dom'

// Polyfills for react-chatbotify which needs browser APIs
Object.defineProperty(global, 'MessageChannel', {
  value: class MessageChannel {
    port1 = {
      postMessage: jest.fn(),
      onmessage: null,
      close: jest.fn(),
      start: jest.fn(),
      addEventListener: jest.fn(),
      removeEventListener: jest.fn(),
      dispatchEvent: jest.fn(),
      onmessageerror: null
    }
    port2 = {
      postMessage: jest.fn(),
      onmessage: null,
      close: jest.fn(),
      start: jest.fn(),
      addEventListener: jest.fn(),
      removeEventListener: jest.fn(),
      dispatchEvent: jest.fn(),
      onmessageerror: null
    }
  }
})

Object.defineProperty(global, 'TextEncoder', {
  value: class TextEncoder {
    encode = jest.fn(() => new Uint8Array())
  }
})

Object.defineProperty(global, 'TextDecoder', {
  value: class TextDecoder {
    decode = jest.fn(() => '')
  }
})