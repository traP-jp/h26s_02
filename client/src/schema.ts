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
  option?: { queryParams?: Record<string, string>; body?: Record<string, any> | FormData }
) => {
  const { userId } = useUserStore()
  const queryParamStr = option?.queryParams
    ? '?' + new URLSearchParams(option.queryParams).toString()
    : ''

  // 送信するデータが FormData かどうかを判定
  const isFormData = option?.body instanceof FormData
  const headers: HeadersInit = { 'X-Forwarded-User': userId! }

  // FormData でない場合のみ JSON 用の Content-Type を設定する
  if (!isFormData) headers['Content-Type'] = 'application/json'
  console.log(`[fetchApi] Request: ${method} ${path}`, { isFormData, body: option?.body })

  const request: RequestInit = {
    method: method,
    headers: headers,
    body:
      option?.body instanceof FormData // もう 1 回直接書かないと型推論がうまくいかない
        ? option.body
        : option?.body
          ? JSON.stringify(option.body)
          : undefined,
  }

  const res = await fetch(`/api${path}${queryParamStr}`, request)
  const text = await res.text()

  if (!res.ok) {
    console.error(`[fetchApi] Error: ${res.status} ${res.statusText}`, text)
    throw new Error(`API Error: ${res.status} ${res.statusText}`)
  } else {
    const result = text ? JSON.parse(text) : {}
    console.log(`[fetchApi] Response: ${res.status}`, result)
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
  newPost: async (body: { image: File; tags: string[] }) => {
    // const validImageTypes = ['image/jpeg', 'image/png']
    // if (!validImageTypes.includes(body.image.type)) {
    //   throw new Error('画像ファイルは png または jpeg 形式である必要があります。')
    // }

    // if (body.tags.length > 10) {
    //   throw new Error('タグは10個より多く指定することはできません。')
    // }

    // for (const tag of body.tags) {
    //   if (tag.includes(',')) {
    //     throw new Error('タグにコンマを含めることはできません。')
    //   }

    //   if (tag.length === 0 || tag.length >= 17) {
    //     throw new Error('タグは1文字以上、16文字以下である必要があります。')
    //   }
    // }

    const formData = new FormData()
    formData.append('image', body.image)
    formData.append('tags', body.tags.length > 0 ? body.tags.join(',') : '')
    return (await fetchApi('POST', '/posts', { body: formData })) as { id: string } // UUID
  },

  getPosts: async (before?: string, limit?: number) => {
    return (await fetchApi('GET', '/posts', {
      queryParams: {
        ...(before ? { before } : {}),
        ...(limit ? { limit: limit.toString() } : {}),
      },
    })) as Post[]
  },

  getPost: async (postId: string) => {
    return (await fetchApi('GET', `/posts/${postId}`)) as Post
  },

  deletePost: async (postId: string) => {
    return await fetchApi('DELETE', `/posts/${postId}`)
  },

  postReaction: async (postId: string, reactionId: number) => {
    return (await fetchApi('POST', `/posts/${postId}/reactions`, {
      body: { id: reactionId },
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
