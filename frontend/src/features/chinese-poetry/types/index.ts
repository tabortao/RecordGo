export interface Poem {
  id: number
  // Core fields from JSON
  cat_cns: string | null
  dynasty_cns: string
  author_cns: string
  title_cns: string
  paragraphs_cns: string[]
  
  // Optional/Nullable fields from JSON
  author_py1?: string
  author_py2?: string
  title_py1?: string
  title_py2?: string
  paragraphs_py1?: string[]
  paragraphs_py2?: string[]
  cat_cnt?: string
  dynasty_cnt?: string
  author_cnt?: string
  title_cnt?: string
  paragraphs_cnt?: string[]
  
  // Custom/App fields
  tags?: string[] // 'primary', 'middle', 'high', etc.
  isCustom?: boolean
  
  // Enriched Content
  translation?: string
  appreciation?: string
  fun_fact?: string
  audio_url?: string
  image_url?: string
}

export interface StudyRecord {
  poemId: number
  lastStudyTime: number // timestamp
  nextReviewTime: number // timestamp
  stage: number // Ebbinghaus stage (0-6)
  timesStudied: number
  isMastered: boolean
  history: number[] // study timestamps
}

export interface DailyStats {
  date: string // YYYY-MM-DD
  studiedCount: number // New + Review
  details: number[] // Poem IDs
}

export interface DictationResult {
  id: string // uuid
  poemId: number
  date: number // timestamp
  accuracy: number // 0-100
  input: string[] // User input lines
  errors: Array<{
    lineIndex: number
    charIndex: number
    expected: string
    actual: string
    type: 'wrong' | 'missing' | 'extra'
  }>
}

export interface PoetryState {
  poems: Poem[]
  studyRecords: Record<number, StudyRecord>
  dailyStats: Record<string, DailyStats>
  collections: number[] // Poem IDs
  dictationRecords: DictationResult[]
}
