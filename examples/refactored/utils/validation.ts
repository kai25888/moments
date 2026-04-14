// utils/validation.ts
// 统一的验证工具函数

import type { ValidationRule, ValidationResult } from '~/types'

/**
 * 验证错误类
 */
export class ValidationError extends Error {
  constructor(
    message: string,
    public field?: string,
    public code?: string
  ) {
    super(message)
    this.name = 'ValidationError'
  }
}

/**
 * 验证规则集合
 */
export const rules = {
  /**
   * 必填验证
   */
  required(message = '此项为必填项'): ValidationRule {
    return {
      validator: (value: unknown) => {
        if (value === undefined || value === null) return false
        if (typeof value === 'string') return value.trim().length > 0
        if (Array.isArray(value)) return value.length > 0
        return true
      },
      message
    }
  },

  /**
   * 最小长度验证
   */
  minLength(min: number, message?: string): ValidationRule {
    return {
      validator: (value: string) => value?.length >= min,
      message: message || `长度不能少于 ${min} 个字符`
    }
  },

  /**
   * 最大长度验证
   */
  maxLength(max: number, message?: string): ValidationRule {
    return {
      validator: (value: string) => !value || value.length <= max,
      message: message || `长度不能超过 ${max} 个字符`
    }
  },

  /**
   * 邮箱格式验证
   */
  email(message = '请输入有效的邮箱地址'): ValidationRule {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
    return {
      validator: (value: string) => !value || emailRegex.test(value),
      message
    }
  },

  /**
   * URL 格式验证
   */
  url(message = '请输入有效的 URL'): ValidationRule {
    return {
      validator: (value: string) => {
        if (!value) return true
        try {
          new URL(value)
          return true
        } catch {
          return false
        }
      },
      message
    }
  },

  /**
   * 正则表达式验证
   */
  pattern(regex: RegExp, message: string): ValidationRule {
    return {
      validator: (value: string) => !value || regex.test(value),
      message
    }
  },

  /**
   * 自定义验证函数
   */
  custom(validator: (value: unknown) => boolean, message: string): ValidationRule {
    return { validator, message }
  }
}

/**
 * 验证单个值
 */
export function validateValue(
  value: unknown,
  validationRules: ValidationRule[]
): ValidationResult {
  for (const rule of validationRules) {
    if (!rule.validator(value)) {
      return {
        valid: false,
        message: rule.message
      }
    }
  }
  return { valid: true }
}

/**
 * 验证对象
 */
export function validateObject<T extends Record<string, unknown>>(
  obj: T,
  schema: Record<keyof T, ValidationRule[]>
): Record<string, string> {
  const errors: Record<string, string> = {}

  for (const [key, rules] of Object.entries(schema)) {
    const value = obj[key as keyof T]
    const result = validateValue(value, rules)
    
    if (!result.valid) {
      errors[key] = result.message
    }
  }

  return errors
}

/**
 * 创建表单验证器
 */
export function createValidator<T extends Record<string, unknown>>(
  schema: Record<keyof T, ValidationRule[]>
) {
  return {
    /**
     * 验证整个表单
     */
    validate(data: T): { valid: boolean; errors: Record<string, string> } {
      const errors = validateObject(data, schema)
      return {
        valid: Object.keys(errors).length === 0,
        errors
      }
    },

    /**
     * 验证单个字段
     */
    validateField<K extends keyof T>(field: K, value: T[K]): string | null {
      const rules = schema[field]
      if (!rules) return null

      const result = validateValue(value, rules)
      return result.valid ? null : result.message
    }
  }
}

// 常用验证模式
export const patterns = {
  // 用户名：字母、数字、下划线，3-20位
  username: /^[a-zA-Z0-9_]{3,20}$/,
  
  // 密码：至少8位，包含字母和数字
  password: /^(?=.*[A-Za-z])(?=.*\d)[A-Za-z\d@$!%*#?&]{8,}$/,
  
  // 手机号（中国大陆）
  phone: /^1[3-9]\d{9}$/,
  
  // 中文姓名
  chineseName: /^[\u4e00-\u9fa5]{2,10}$/,
  
  // 标签：不允许特殊字符
  tag: /^[\u4e00-\u9fa5a-zA-Z0-9_-]+$/
}

// 预设验证器
export const validators = {
  /**
   * 用户名验证
   */
  username: [
    rules.required('请输入用户名'),
    rules.pattern(patterns.username, '用户名只能包含字母、数字、下划线，长度3-20位')
  ],

  /**
   * 密码验证
   */
  password: [
    rules.required('请输入密码'),
    rules.pattern(patterns.password, '密码至少8位，需包含字母和数字')
  ],

  /**
   * 邮箱验证
   */
  email: [
    rules.required('请输入邮箱'),
    rules.email()
  ],

  /**
   * 手机号验证
   */
  phone: [
    rules.pattern(patterns.phone, '请输入有效的手机号')
  ],

  /**
   * Memo 内容验证
   */
  memoContent: [
    rules.required('请输入内容'),
    rules.maxLength(10000, '内容不能超过10000字')
  ],

  /**
   * 标签验证
   */
  tag: [
    rules.required('请输入标签'),
    rules.maxLength(20, '标签不能超过20个字符'),
    rules.pattern(patterns.tag, '标签只能包含中文、字母、数字、下划线和横线')
  ]
}
