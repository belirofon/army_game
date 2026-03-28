import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import { defineComponent, h, computed } from 'vue'

const StatCard = defineComponent({
  props: {
    name: { type: String, required: true },
    value: { type: Number, required: true },
    max: { type: Number, default: 100 },
    description: { type: String, default: '' },
  },
  setup(props) {
    const percentage = computed(() => Math.round((props.value / props.max) * 100))
    const statusClass = computed(() => {
      if (percentage.value < 30) return 'stat-critical'
      if (percentage.value < 50) return 'stat-warning'
      return 'stat-normal'
    })
    return () => h('div', { class: 'stat-card' }, [
      h('div', { class: 'stat-header' }, [
        h('span', { class: 'stat-name' }, props.name),
        h('span', { class: `stat-value ${statusClass.value}` }, `${props.value}/${props.max}`),
      ]),
      h('div', { class: 'stat-bar-bg' }, [
        h('div', { 
          class: `stat-bar-fill ${statusClass.value}`, 
          style: { width: `${percentage.value}%` } 
        }),
      ]),
      props.description ? h('p', { class: 'stat-description' }, props.description) : null,
    ])
  },
})

describe('StatCard Component', () => {
  it('renders stat name', () => {
    const wrapper = mount(StatCard, { props: { name: 'STR', value: 50 } })
    expect(wrapper.find('.stat-name').text()).toBe('STR')
  })

  it('renders stat value and max', () => {
    const wrapper = mount(StatCard, { props: { name: 'STR', value: 75, max: 100 } })
    expect(wrapper.find('.stat-value').text()).toBe('75/100')
  })

  it('applies stat-normal class when value >= 50', () => {
    const wrapper = mount(StatCard, { props: { name: 'STR', value: 60 } })
    expect(wrapper.find('.stat-normal').exists()).toBe(true)
  })

  it('applies stat-warning class when value < 50', () => {
    const wrapper = mount(StatCard, { props: { name: 'STR', value: 40 } })
    expect(wrapper.find('.stat-warning').exists()).toBe(true)
  })

  it('applies stat-critical class when value < 30', () => {
    const wrapper = mount(StatCard, { props: { name: 'STR', value: 20 } })
    expect(wrapper.find('.stat-critical').exists()).toBe(true)
  })

  it('renders description when provided', () => {
    const wrapper = mount(StatCard, { props: { name: 'STR', value: 50, description: 'Strength' } })
    expect(wrapper.find('.stat-description').text()).toBe('Strength')
  })

  it('does not render description when not provided', () => {
    const wrapper = mount(StatCard, { props: { name: 'STR', value: 50 } })
    expect(wrapper.find('.stat-description').exists()).toBe(false)
  })

  it('uses custom max value', () => {
    const wrapper = mount(StatCard, { props: { name: 'MOR', value: 50, max: 200 } })
    expect(wrapper.find('.stat-value').text()).toBe('50/200')
  })
})
