import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import { defineComponent, h } from 'vue'

const Spinner = defineComponent({
  props: {
    size: { type: String, default: 'md' },
  },
  setup(props) {
    return () => h('div', { class: `spinner-wrapper spinner-${props.size}` }, [
      h('div', { class: 'spinner-circle' }),
    ])
  },
})

describe('Spinner Component', () => {
  it('renders spinner wrapper', () => {
    const wrapper = mount(Spinner)
    expect(wrapper.find('.spinner-wrapper').exists()).toBe(true)
  })

  it('renders spinner circle', () => {
    const wrapper = mount(Spinner)
    expect(wrapper.find('.spinner-circle').exists()).toBe(true)
  })

  it('applies sm size class', () => {
    const wrapper = mount(Spinner, { props: { size: 'sm' } })
    expect(wrapper.find('.spinner-sm').exists()).toBe(true)
  })

  it('applies md size class by default', () => {
    const wrapper = mount(Spinner)
    expect(wrapper.find('.spinner-md').exists()).toBe(true)
  })

  it('applies lg size class', () => {
    const wrapper = mount(Spinner, { props: { size: 'lg' } })
    expect(wrapper.find('.spinner-lg').exists()).toBe(true)
  })
})
