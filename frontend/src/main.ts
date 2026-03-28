import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { ApolloClient, InMemoryCache, createHttpLink } from '@apollo/client/core'
import { DefaultApolloClient } from '@vue/apollo-composable'
import App from './App.vue'
import router from './router'
import './assets/styles/main.css'

const httpLink = createHttpLink({
  uri: '/graphql',
})

const apolloClient = new ApolloClient({
  link: httpLink,
  cache: new InMemoryCache(),
  defaultOptions: {
    watchQuery: { fetchPolicy: 'network-only' },
    query: { fetchPolicy: 'network-only' },
  },
})

const app = createApp(App)
const pinia = createPinia()

app.provide(DefaultApolloClient, apolloClient)
app.use(pinia)
app.use(router)

app.mount('#app')