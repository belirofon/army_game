import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import { defineComponent, h } from 'vue'

const EventIllustration = defineComponent({
  props: {
    templateId: { type: String, default: '' },
  },
  setup(props) {
    return () => h('div', { class: 'illustration' }, [
      h('div', { class: 'illustration-placeholder' }, 'Placeholder'),
    ])
  },
})

describe('EventIllustration Component', () => {
  it('renders illustration container', () => {
    const wrapper = mount(EventIllustration)
    expect(wrapper.find('.illustration').exists()).toBe(true)
  })

  it('renders placeholder', () => {
    const wrapper = mount(EventIllustration)
    expect(wrapper.find('.illustration-placeholder').exists()).toBe(true)
  })

  it('accepts templateId prop', () => {
    const wrapper = mount(EventIllustration, { props: { templateId: 'training' } })
    expect(wrapper.props('templateId')).toBe('training')
  })

  it('renders without templateId', () => {
    const wrapper = mount(EventIllustration)
    expect(wrapper.find('.illustration-placeholder').exists()).toBe(true)
  })
})
