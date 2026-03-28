import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import { defineComponent, h } from 'vue'

const Card = defineComponent({
  props: {
    variant: { type: String, default: 'default' },
  },
  setup(props, { slots }) {
    return () => h('div', { class: `card ${props.variant === 'elevated' ? 'card-elevated' : ''}` }, slots.default?.())
  },
})

describe('Card Component', () => {
  it('renders card element', () => {
    const wrapper = mount(Card)
    expect(wrapper.find('.card').exists()).toBe(true)
  })

  it('renders default variant class', () => {
    const wrapper = mount(Card)
    expect(wrapper.find('.card').classes()).toContain('card')
  })

  it('applies default variant when not specified', () => {
    const wrapper = mount(Card)
    expect(wrapper.find('.card-elevated').exists()).toBe(false)
  })

  it('applies elevated variant class', () => {
    const wrapper = mount(Card, { props: { variant: 'elevated' } })
    expect(wrapper.find('.card-elevated').exists()).toBe(true)
  })

  it('renders slot content', () => {
    const wrapper = mount(Card, { slots: { default: 'Card Content' } })
    expect(wrapper.text()).toBe('Card Content')
  })
})
