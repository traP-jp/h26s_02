/* eslint-disable
  @typescript-eslint/no-unsafe-return,
  @typescript-eslint/no-explicit-any,
  @typescript-eslint/no-unsafe-assignment
*/
import { useUserStore } from '@/store'

type HttpMethod = 'GET' | 'POST' | 'PUT' | 'DELETE' | 'PATCH'

const fetchApi = async (
  method: HttpMethod,
  path: string,
  option?: { queryParams?: Record<string, string>; body?: Record<string, any> }
) => {
  const { userId } = useUserStore()
  const queryParamStr = option?.queryParams
    ? '?' + new URLSearchParams(option.queryParams).toString()
    : ''

  const request: RequestInit = {
    method: method,
    headers: { 'X-Forwarded-User': userId!, 'Content-Type': 'application/json' },
    body: option?.body ? JSON.stringify(option.body) : undefined,
  }

  const res = await fetch(`/api${path}${queryParamStr}`, request)
  const text = await res.text()
  if (!res.ok) {
    throw new Error(`API Error: ${res.status} ${res.statusText}`)
  } else {
    const result = text ? JSON.parse(text) : {}
    // console.log(result)
    return result
  }
}

export type Reaction = {
  id: number // 暫定 5 以下
  count: number
}

export type TagInfo = {
  name: string
  count: number
}

export type Post = {
  id: string // UUID
  userName: string // traQ ID
  tags: string[]
  imageUrl: string
  reactions: Reaction[]
  createdAt: string // ISO 8601
}

export const api = {
  newPost: async (body: { image: string; tags: string[] }) => {
    return await fetchApi('POST', '/posts', { body })
  },

  getPosts: async () => {
    return (await fetchApi('GET', '/posts')) as Post[]
  },

  getPost: async (postId: string) => {
    return (await fetchApi('GET', `/posts/${postId}`)) as Post
  },

  deletePost: async (postId: string) => {
    return await fetchApi('DELETE', `/posts/${postId}`)
  },

  postReaction: async (postId: string, reactionId: number) => {
    return (await fetchApi('POST', `/posts/${postId}/reactions`, {
      body: { reactionId },
    })) as Reaction[]
  },

  deleteReaction: async (postId: string, reactionId: number) => {
    return await fetchApi('DELETE', `/posts/${postId}/reactions/${reactionId}`)
  },

  getTags: async () => {
    return (await fetchApi('GET', '/tags')) as TagInfo[]
  },

  getTaggedPosts: async (tag: string) => {
    return (await fetchApi('GET', '/posts', { queryParams: { tag } })) as Post[]
  },

  getMe: async () => {
    return (await fetchApi('GET', '/users/me')) as { userName: string }
  },

  getUserPosts: async (userName: string) => {
    return (await fetchApi('GET', `/users/${userName}/posts`)) as Post[]
  },
}
