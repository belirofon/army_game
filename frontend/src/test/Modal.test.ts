import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import { defineComponent, h } from 'vue'

const Modal = defineComponent({
  props: {
    show: { type: Boolean, required: true },
    title: { type: String, default: '' },
  },
  emits: ['close'],
  setup(props, { emit, slots }) {
    return () => props.show ? h('div', { class: 'modal-overlay', onClick: () => emit('close') }, [
      h('div', { class: 'modal-container', onClick: (e: MouseEvent) => e.stopPropagation() }, [
        props.title ? h('div', { class: 'modal-header' }, [
          h('h2', { class: 'modal-title' }, props.title),
          h('button', { class: 'modal-close', onClick: () => emit('close') }, '×'),
        ]) : null,
        h('div', { class: 'modal-content' }, slots.default?.()),
      ]),
    ]) : null
  },
})

describe('Modal Component', () => {
  it('does not render when show is false', () => {
    const wrapper = mount(Modal, { props: { show: false } })
    expect(wrapper.find('.modal-overlay').exists()).toBe(false)
  })

  it('renders when show is true', () => {
    const wrapper = mount(Modal, { props: { show: true } })
    expect(wrapper.find('.modal-overlay').exists()).toBe(true)
  })

  it('renders title when provided', () => {
    const wrapper = mount(Modal, { props: { show: true, title: 'Test Title' } })
    expect(wrapper.find('.modal-title').text()).toBe('Test Title')
  })

  it('does not render title when not provided', () => {
    const wrapper = mount(Modal, { props: { show: true } })
    expect(wrapper.find('.modal-header').exists()).toBe(false)
  })

  it('renders content slot', () => {
    const wrapper = mount(Modal, { 
      props: { show: true }, 
      slots: { default: () => h('span', 'Test Content') } 
    })
    expect(wrapper.find('.modal-content').text()).toBe('Test Content')
  })

  it('emits close when backdrop clicked', async () => {
    const wrapper = mount(Modal, { props: { show: true } })
    await wrapper.find('.modal-overlay').trigger('click')
    expect(wrapper.emitted('close')).toBeTruthy()
  })

  it('does not emit close when content clicked', async () => {
    const wrapper = mount(Modal, { props: { show: true } })
    await wrapper.find('.modal-container').trigger('click')
    expect(wrapper.emitted('close')).toBeFalsy()
  })

  it('emits close when close button clicked', async () => {
    const wrapper = mount(Modal, { props: { show: true, title: 'Test' } })
    await wrapper.find('.modal-close').trigger('click')
    expect(wrapper.emitted('close')).toBeTruthy()
  })
})
