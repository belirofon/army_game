import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import { defineComponent, h } from 'vue'

const Toast = defineComponent({
  props: {
    message: { type: String, required: true },
    type: { type: String, default: 'success' },
  },
  emits: ['close'],
  setup(props, { emit }) {
    return () => h('div', { class: `toast toast-${props.type}` }, [
      h('div', { class: 'toast-icon' }, props.type),
      h('span', { class: 'toast-message' }, props.message),
      h('button', { class: 'toast-close', onClick: () => emit('close') }, '×'),
    ])
  },
})

describe('Toast Component', () => {
  it('renders message', () => {
    const wrapper = mount(Toast, { props: { message: 'Test message' } })
    expect(wrapper.find('.toast-message').text()).toBe('Test message')
  })

  it('applies success type class by default', () => {
    const wrapper = mount(Toast, { props: { message: 'Test' } })
    expect(wrapper.find('.toast-success').exists()).toBe(true)
  })

  it('applies error type class', () => {
    const wrapper = mount(Toast, { props: { message: 'Test', type: 'error' } })
    expect(wrapper.find('.toast-error').exists()).toBe(true)
  })

  it('applies warning type class', () => {
    const wrapper = mount(Toast, { props: { message: 'Test', type: 'warning' } })
    expect(wrapper.find('.toast-warning').exists()).toBe(true)
  })

  it('renders success icon', () => {
    const wrapper = mount(Toast, { props: { message: 'Test', type: 'success' } })
    expect(wrapper.find('.toast-icon').text()).toBe('success')
  })

  it('renders error icon', () => {
    const wrapper = mount(Toast, { props: { message: 'Test', type: 'error' } })
    expect(wrapper.find('.toast-icon').text()).toBe('error')
  })

  it('renders warning icon', () => {
    const wrapper = mount(Toast, { props: { message: 'Test', type: 'warning' } })
    expect(wrapper.find('.toast-icon').text()).toBe('warning')
  })

  it('renders close button', () => {
    const wrapper = mount(Toast, { props: { message: 'Test' } })
    expect(wrapper.find('.toast-close').exists()).toBe(true)
  })

  it('emits close when close button clicked', async () => {
    const wrapper = mount(Toast, { props: { message: 'Test' } })
    await wrapper.find('.toast-close').trigger('click')
    expect(wrapper.emitted('close')).toBeTruthy()
  })
})
