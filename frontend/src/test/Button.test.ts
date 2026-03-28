import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import { defineComponent, h } from 'vue'

const Button = defineComponent({
  props: {
    text: { type: String, required: true },
    variant: { type: String, default: 'primary' },
    disabled: { type: Boolean, default: false },
    loading: { type: Boolean, default: false },
  },
  emits: ['click'],
  setup(props, { emit }) {
    return () => h('button', {
      class: [
        'btn',
        `btn-${props.variant}`,
        { 'btn-disabled': props.disabled || props.loading },
      ],
      disabled: props.disabled || props.loading,
      onClick: () => emit('click'),
    }, [
      props.loading ? h('span', { class: 'spinner' }) : props.text,
    ])
  },
})

describe('Button Component', () => {
  it('renders with text', () => {
    const wrapper = mount(Button, {
      props: { text: 'Click me' },
    })
    
    expect(wrapper.text()).toBe('Click me')
    expect(wrapper.find('button').exists()).toBe(true)
  })

  it('applies primary variant class by default', () => {
    const wrapper = mount(Button, {
      props: { text: 'Primary' },
    })
    
    expect(wrapper.find('.btn-primary').exists()).toBe(true)
  })

  it('applies secondary variant class', () => {
    const wrapper = mount(Button, {
      props: { text: 'Secondary', variant: 'secondary' },
    })
    
    expect(wrapper.find('.btn-secondary').exists()).toBe(true)
  })

  it('applies choice variant class', () => {
    const wrapper = mount(Button, {
      props: { text: 'Choice', variant: 'choice' },
    })
    
    expect(wrapper.find('.btn-choice').exists()).toBe(true)
  })

  it('emits click event when clicked', () => {
    const wrapper = mount(Button, {
      props: { text: 'Click me' },
    })
    
    wrapper.find('button').trigger('click')
    
    expect(wrapper.emitted('click')).toBeTruthy()
    expect(wrapper.emitted('click')?.length).toBe(1)
  })

  it('does not emit click when disabled', () => {
    const wrapper = mount(Button, {
      props: { text: 'Disabled', disabled: true },
    })
    
    wrapper.find('button').trigger('click')
    
    expect(wrapper.emitted('click')).toBeFalsy()
  })

  it('does not emit click when loading', () => {
    const wrapper = mount(Button, {
      props: { text: 'Loading', loading: true },
    })
    
    wrapper.find('button').trigger('click')
    
    expect(wrapper.emitted('click')).toBeFalsy()
  })

  it('shows spinner when loading', () => {
    const wrapper = mount(Button, {
      props: { text: 'Loading', loading: true },
    })
    
    expect(wrapper.find('.spinner').exists()).toBe(true)
    expect(wrapper.text()).not.toBe('Loading')
  })

  it('does not show spinner when not loading', () => {
    const wrapper = mount(Button, {
      props: { text: 'Not loading', loading: false },
    })
    
    expect(wrapper.find('.spinner').exists()).toBe(false)
  })

  it('applies disabled class when disabled', () => {
    const wrapper = mount(Button, {
      props: { text: 'Disabled', disabled: true },
    })
    
    expect(wrapper.find('.btn-disabled').exists()).toBe(true)
  })

  it('applies disabled class when loading', () => {
    const wrapper = mount(Button, {
      props: { text: 'Loading', loading: true },
    })
    
    expect(wrapper.find('.btn-disabled').exists()).toBe(true)
  })

  it('button has disabled attribute when disabled', () => {
    const wrapper = mount(Button, {
      props: { text: 'Disabled', disabled: true },
    })
    
    expect(wrapper.find('button').attributes('disabled')).toBeDefined()
  })

  it('button has disabled attribute when loading', () => {
    const wrapper = mount(Button, {
      props: { text: 'Loading', loading: true },
    })
    
    expect(wrapper.find('button').attributes('disabled')).toBeDefined()
  })

  it('button does not have disabled attribute when enabled', () => {
    const wrapper = mount(Button, {
      props: { text: 'Enabled', disabled: false, loading: false },
    })
    
    expect(wrapper.find('button').attributes('disabled')).toBeUndefined()
  })
})
