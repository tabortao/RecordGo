// 拼音数据字典

// 声母表及其对应图标（用于转盘展示）
export const initialsData: Record<string, { icon: string; label: string }> = {
  b: { icon: "🍍", label: "菠萝" },
  p: { icon: "🍇", label: "葡萄" },
  m: { icon: "🐎", label: "马" },
  f: { icon: "🪓", label: "斧头" },
  d: { icon: "🥁", label: "鼓" },
  t: { icon: "🍑", label: "桃子" },
  n: { icon: "🎃", label: "南瓜" },
  l: { icon: "🍐", label: "梨" },
  g: { icon: "🕊️", label: "鸽子" },
  k: { icon: "🐸", label: "蝌蚪" },
  h: { icon: "🐵", label: "猴子" },
  j: { icon: "🐔", label: "公鸡" },
  q: { icon: "🎈", label: "气球" },
  x: { icon: "🍉", label: "西瓜" },
  zh: { icon: "🕷️", label: "蜘蛛" },
  ch: { icon: "📏", label: "尺子" },
  sh: { icon: "🦁", label: "狮子" },
  r: { icon: "☀️", label: "日" },
  z: { icon: "✍️", label: "写字" },
  c: { icon: "🦔", label: "刺猬" },
  s: { icon: "🐍", label: "蛇" },
  y: { icon: "🦆", label: "鸭子" },
  w: { icon: "🐌", label: "蜗牛" },
};

export const initials = Object.keys(initialsData);

// 韵母表
export const finals = [
  "a",
  "o",
  "e",
  "i",
  "u",
  "ü",
  "ai",
  "ei",
  "ui",
  "ao",
  "ou",
  "iu",
  "ie",
  "üe",
  "er",
  "an",
  "en",
  "in",
  "un",
  "ün",
  "ang",
  "eng",
  "ing",
  "ong",
];

// 声调
export const tones = [1, 2, 3, 4];

// 拼音组合字典
// 格式: key = "initial+final", value = { 1: { word: "...", pinyin: "..." }, ... }
export const pinyinDict: Record<
  string,
  Record<number, { word: string; pinyin: string; example?: string }>
> = {
  // b
  ba: {
    1: { word: "八", pinyin: "bā", example: "数字八" },
    2: { word: "拔", pinyin: "bá", example: "拔萝卜" },
    3: { word: "靶", pinyin: "bǎ", example: "靶心" },
    4: { word: "爸", pinyin: "bà", example: "爸爸" },
  },
  bo: {
    1: { word: "波", pinyin: "bō", example: "波浪" },
    2: { word: "伯", pinyin: "bó", example: "伯伯" },
    3: { word: "跛", pinyin: "bǒ", example: "跛脚" },
    4: { word: "薄", pinyin: "bò", example: "薄荷" },
  },
  bi: {
    1: { word: "笔", pinyin: "bī", example: "铅笔" }, // 这里的声调其实多为3声，但在游戏中可以作为展示
    2: { word: "鼻", pinyin: "bí", example: "鼻子" },
    3: { word: "笔", pinyin: "bǐ", example: "写字" },
    4: { word: "币", pinyin: "bì", example: "硬币" },
  },
  bu: {
    1: { word: "逋", pinyin: "bū", example: "逃脱" },
    2: { word: "补", pinyin: "bú", example: "补充" }, // 变调，这里简化
    3: { word: "补", pinyin: "bǔ", example: "缝补" },
    4: { word: "布", pinyin: "bù", example: "布料" },
  },
  bai: {
    1: { word: "掰", pinyin: "bāi", example: "掰开" },
    2: { word: "白", pinyin: "bái", example: "白色" },
    3: { word: "百", pinyin: "bǎi", example: "一百" },
    4: { word: "拜", pinyin: "bài", example: "拜年" },
  },

  // p
  pa: {
    1: { word: "趴", pinyin: "pā", example: "趴下" },
    2: { word: "爬", pinyin: "pá", example: "爬山" },
    3: { word: "怕", pinyin: "pǎ", example: "害怕" }, // 口语中怕通常是4声，这里为了完整性
    4: { word: "怕", pinyin: "pà", example: "害怕" },
  },
  po: {
    1: { word: "坡", pinyin: "pō", example: "山坡" },
    2: { word: "婆", pinyin: "pó", example: "外婆" },
    3: { word: "叵", pinyin: "pǒ", example: "居心叵测" },
    4: { word: "破", pinyin: "pò", example: "破坏" },
  },

  // m
  ma: {
    1: { word: "妈", pinyin: "mā", example: "妈妈" },
    2: { word: "麻", pinyin: "má", example: "芝麻" },
    3: { word: "马", pinyin: "mǎ", example: "骑马" },
    4: { word: "骂", pinyin: "mà", example: "责骂" },
  },
  mi: {
    1: { word: "咪", pinyin: "mī", example: "猫咪" },
    2: { word: "迷", pinyin: "mí", example: "迷路" },
    3: { word: "米", pinyin: "mǐ", example: "大米" },
    4: { word: "密", pinyin: "mì", example: "秘密" },
  },
  mu: {
    1: { word: "摸", pinyin: "mū", example: "摸" }, // 特殊音
    2: { word: "模", pinyin: "mú", example: "模样" },
    3: { word: "母", pinyin: "mǔ", example: "母亲" },
    4: { word: "木", pinyin: "mù", example: "树木" },
  },

  // f
  fa: {
    1: { word: "发", pinyin: "fā", example: "出发" },
    2: { word: "罚", pinyin: "fá", example: "惩罚" },
    3: { word: "法", pinyin: "fǎ", example: "办法" },
    4: { word: "发", pinyin: "fà", example: "头发" },
  },
  fu: {
    1: { word: "夫", pinyin: "fū", example: "大夫" },
    2: { word: "福", pinyin: "fú", example: "福气" },
    3: { word: "府", pinyin: "fǔ", example: "政府" },
    4: { word: "父", pinyin: "fù", example: "父亲" },
  },

  // d
  da: {
    1: { word: "搭", pinyin: "dā", example: "搭积木" },
    2: { word: "答", pinyin: "dá", example: "回答" },
    3: { word: "打", pinyin: "dǎ", example: "打球" },
    4: { word: "大", pinyin: "dà", example: "大象" },
  },
  de: {
    1: { word: "得", pinyin: "dē", example: "得" },
    2: { word: "德", pinyin: "dé", example: "品德" },
    3: { word: "得", pinyin: "dě", example: "得" }, // 很少用
    4: { word: "特", pinyin: "tè", example: "特别" }, // 借用近似
  },

  // t
  tu: {
    1: { word: "秃", pinyin: "tū", example: "秃头" },
    2: { word: "图", pinyin: "tú", example: "图画" },
    3: { word: "土", pinyin: "tǔ", example: "土地" },
    4: { word: "兔", pinyin: "tù", example: "兔子" },
  },
  ta: {
    1: { word: "他", pinyin: "tā", example: "他们" },
    2: { word: "塔", pinyin: "tá", example: "宝塔" },
    3: { word: "塔", pinyin: "tǎ", example: "塔" },
    4: { word: "踏", pinyin: "tà", example: "踏步" },
  },

  // n
  na: {
    1: { word: "娜", pinyin: "nā", example: "安娜" },
    2: { word: "拿", pinyin: "ná", example: "拿走" },
    3: { word: "哪", pinyin: "nǎ", example: "哪里" },
    4: { word: "那", pinyin: "nà", example: "那里" },
  },
  ni: {
    1: { word: "妮", pinyin: "nī", example: "妮妮" },
    2: { word: "泥", pinyin: "ní", example: "泥土" },
    3: { word: "你", pinyin: "nǐ", example: "你好" },
    4: { word: "逆", pinyin: "nì", example: "逆风" },
  },
  nü: {
    3: { word: "女", pinyin: "nǚ", example: "女孩" },
    4: { word: "怒", pinyin: "nù", example: "发怒" }, // 近似
  },

  // l
  la: {
    1: { word: "拉", pinyin: "lā", example: "拉手" },
    2: { word: "啦", pinyin: "lá", example: "啦啦队" },
    3: { word: "喇", pinyin: "lǎ", example: "喇叭" },
    4: { word: "辣", pinyin: "là", example: "辣椒" },
  },
  le: {
    4: { word: "乐", pinyin: "lè", example: "快乐" },
  },
  li: {
    1: { word: "哩", pinyin: "lī", example: "哩" },
    2: { word: "梨", pinyin: "lí", example: "鸭梨" },
    3: { word: "李", pinyin: "lǐ", example: "李子" },
    4: { word: "力", pinyin: "lì", example: "力量" },
  },
  lü: {
    3: { word: "旅", pinyin: "lǚ", example: "旅游" },
    4: { word: "绿", pinyin: "lǜ", example: "绿色" },
  },

  // g
  ge: {
    1: { word: "歌", pinyin: "gē", example: "唱歌" },
    2: { word: "格", pinyin: "gé", example: "格子" },
    3: { word: "葛", pinyin: "gě", example: "诸葛亮" },
    4: { word: "个", pinyin: "gè", example: "一个" },
  },
  gu: {
    1: { word: "姑", pinyin: "gū", example: "姑娘" },
    2: { word: "骨", pinyin: "gú", example: "骨头" },
    3: { word: "鼓", pinyin: "gǔ", example: "打鼓" },
    4: { word: "故", pinyin: "gù", example: "故事" },
  },

  // k
  ka: {
    1: { word: "咖", pinyin: "kā", example: "咖啡" },
    3: { word: "卡", pinyin: "kǎ", example: "卡片" },
  },
  ke: {
    1: { word: "科", pinyin: "kē", example: "科学" },
    2: { word: "壳", pinyin: "ké", example: "贝壳" },
    3: { word: "可", pinyin: "kě", example: "可以" },
    4: { word: "课", pinyin: "kè", example: "上课" },
  },

  // h
  he: {
    1: { word: "喝", pinyin: "hē", example: "喝水" },
    2: { word: "河", pinyin: "hé", example: "小河" },
    4: { word: "贺", pinyin: "hè", example: "祝贺" },
  },
  hu: {
    1: { word: "呼", pinyin: "hū", example: "呼吸" },
    2: { word: "壶", pinyin: "hú", example: "水壶" },
    3: { word: "虎", pinyin: "hǔ", example: "老虎" },
    4: { word: "护", pinyin: "hù", example: "保护" },
  },

  // j
  ji: {
    1: { word: "鸡", pinyin: "jī", example: "公鸡" },
    2: { word: "急", pinyin: "jí", example: "着急" },
    3: { word: "几", pinyin: "jǐ", example: "几个" },
    4: { word: "记", pinyin: "jì", example: "日记" },
  },
  ju: {
    1: { word: "居", pinyin: "jū", example: "居住" },
    2: { word: "局", pinyin: "jú", example: "邮局" },
    3: { word: "举", pinyin: "jǔ", example: "举手" },
    4: { word: "巨", pinyin: "jù", example: "巨大" },
  },

  // q
  qi: {
    1: { word: "七", pinyin: "qī", example: "数字七" },
    2: { word: "齐", pinyin: "qí", example: "整齐" },
    3: { word: "起", pinyin: "qǐ", example: "起床" },
    4: { word: "气", pinyin: "qì", example: "生气" },
  },
  qu: {
    1: { word: "区", pinyin: "qū", example: "地区" },
    2: { word: "渠", pinyin: "qú", example: "水渠" },
    3: { word: "取", pinyin: "qǔ", example: "取出" },
    4: { word: "去", pinyin: "qù", example: "回去" },
  },

  // x
  xi: {
    1: { word: "西", pinyin: "xī", example: "西瓜" },
    2: { word: "习", pinyin: "xí", example: "学习" },
    3: { word: "洗", pinyin: "xǐ", example: "洗手" },
    4: { word: "戏", pinyin: "xì", example: "游戏" },
  },
  xu: {
    1: { word: "需", pinyin: "xū", example: "需要" },
    2: { word: "徐", pinyin: "xú", example: "徐徐" },
    3: { word: "许", pinyin: "xǔ", example: "许多" },
    4: { word: "序", pinyin: "xù", example: "顺序" },
  },

  // zh
  zhu: {
    1: { word: "猪", pinyin: "zhū", example: "小猪" },
    2: { word: "竹", pinyin: "zhú", example: "竹子" },
    3: { word: "主", pinyin: "zhǔ", example: "主人" },
    4: { word: "住", pinyin: "zhù", example: "住房" },
  },
  zha: {
    1: { word: "扎", pinyin: "zhā", example: "扎针" },
    2: { word: "炸", pinyin: "zhá", example: "油炸" },
    3: { word: "眨", pinyin: "zhǎ", example: "眨眼" },
    4: { word: "炸", pinyin: "zhà", example: "爆炸" },
  },

  // ch
  chi: {
    1: { word: "吃", pinyin: "chī", example: "吃饭" },
    2: { word: "迟", pinyin: "chí", example: "迟到" },
    3: { word: "尺", pinyin: "chǐ", example: "尺子" },
    4: { word: "翅", pinyin: "chì", example: "翅膀" },
  },
  cha: {
    1: { word: "插", pinyin: "chā", example: "插队" },
    2: { word: "茶", pinyin: "chá", example: "喝茶" },
    3: { word: "叉", pinyin: "chǎ", example: "叉子" },
    4: { word: "差", pinyin: "chà", example: "差劲" },
  },

  // sh
  shu: {
    1: { word: "书", pinyin: "shū", example: "看书" },
    2: { word: "熟", pinyin: "shú", example: "成熟" },
    3: { word: "鼠", pinyin: "shǔ", example: "老鼠" },
    4: { word: "树", pinyin: "shù", example: "大树" },
  },
  sha: {
    1: { word: "沙", pinyin: "shā", example: "沙子" },
    3: { word: "傻", pinyin: "shǎ", example: "傻瓜" },
    4: { word: "厦", pinyin: "shà", example: "大厦" },
  },

  // r
  ri: { 4: { word: "日", pinyin: "rì", example: "日出" } },
  re: { 4: { word: "热", pinyin: "rè", example: "好热" } },

  // z
  za: {
    1: { word: "扎", pinyin: "zā", example: "扎辫子" },
    2: { word: "杂", pinyin: "zá", example: "杂技" },
  },
  zu: {
    1: { word: "租", pinyin: "zū", example: "出租" },
    2: { word: "足", pinyin: "zú", example: "足球" },
    3: { word: "组", pinyin: "zǔ", example: "小组" },
  },

  // c
  ci: {
    1: { word: "刺", pinyin: "cī", example: "刺猬" },
    2: { word: "词", pinyin: "cí", example: "词语" },
    3: { word: "此", pinyin: "cǐ", example: "此时" },
    4: { word: "次", pinyin: "cì", example: "次数" },
  },
  ca: {
    1: { word: "擦", pinyin: "cā", example: "擦黑板" },
  },

  // s
  si: {
    1: { word: "丝", pinyin: "sī", example: "丝瓜" },
    3: { word: "死", pinyin: "sǐ", example: "枯死" },
    4: { word: "四", pinyin: "sì", example: "数字四" },
  },
  sa: {
    1: { word: "撒", pinyin: "sā", example: "撒娇" },
    3: { word: "洒", pinyin: "sǎ", example: "洒水" },
  },

  // y
  ya: {
    1: { word: "鸭", pinyin: "yā", example: "鸭子" },
    2: { word: "牙", pinyin: "yá", example: "牙齿" },
    3: { word: "雅", pinyin: "yǎ", example: "优雅" },
    4: { word: "亚", pinyin: "yà", example: "亚洲" },
  },
  ye: {
    1: { word: "耶", pinyin: "yē", example: "耶" },
    2: { word: "爷", pinyin: "yé", example: "爷爷" },
    3: { word: "野", pinyin: "yě", example: "野外" },
    4: { word: "夜", pinyin: "yè", example: "夜晚" },
  },
  yi: {
    1: { word: "衣", pinyin: "yī", example: "衣服" },
    2: { word: "姨", pinyin: "yí", example: "阿姨" },
    3: { word: "椅", pinyin: "yǐ", example: "椅子" },
    4: { word: "意", pinyin: "yì", example: "意思" },
  },
  yu: {
    1: { word: "迂", pinyin: "yū", example: "迂回" },
    2: { word: "鱼", pinyin: "yú", example: "小鱼" },
    3: { word: "雨", pinyin: "yǔ", example: "下雨" },
    4: { word: "玉", pinyin: "yù", example: "玉米" },
  },
  yue: {
    1: { word: "约", pinyin: "yuē", example: "约定" },
    4: { word: "月", pinyin: "yuè", example: "月亮" },
  },
  yun: {
    1: { word: "晕", pinyin: "yūn", example: "晕倒" },
    2: { word: "云", pinyin: "yún", example: "白云" },
    3: { word: "允", pinyin: "yǔn", example: "允许" },
    4: { word: "运", pinyin: "yùn", example: "运动" },
  },
  yuan: {
    1: { word: "冤", pinyin: "yuān", example: "冤枉" },
    2: { word: "圆", pinyin: "yuán", example: "圆形" },
    3: { word: "远", pinyin: "yuǎn", example: "远方" },
    4: { word: "院", pinyin: "yuàn", example: "院子" },
  },
  yin: {
    1: { word: "音", pinyin: "yīn", example: "音乐" },
    2: { word: "银", pinyin: "yín", example: "银行" },
    3: { word: "引", pinyin: "yǐn", example: "引导" },
    4: { word: "印", pinyin: "yìn", example: "脚印" },
  },
  ying: {
    1: { word: "英", pinyin: "yīng", example: "英雄" },
    2: { word: "赢", pinyin: "yíng", example: "赢了" },
    3: { word: "影", pinyin: "yǐng", example: "影子" },
    4: { word: "硬", pinyin: "yìng", example: "坚硬" },
  },
  yang: {
    1: { word: "央", pinyin: "yāng", example: "中央" },
    2: { word: "羊", pinyin: "yáng", example: "山羊" },
    3: { word: "养", pinyin: "yǎng", example: "养育" },
    4: { word: "样", pinyin: "yàng", example: "样子" },
  },
  yong: {
    1: { word: "拥", pinyin: "yōng", example: "拥抱" },
    3: { word: "勇", pinyin: "yǒng", example: "勇敢" },
    4: { word: "用", pinyin: "yòng", example: "有用" },
  },
  you: {
    1: { word: "优", pinyin: "yōu", example: "优秀" },
    2: { word: "游", pinyin: "yóu", example: "游泳" },
    3: { word: "有", pinyin: "yǒu", example: "拥有" },
    4: { word: "右", pinyin: "yòu", example: "右手" },
  },

  // w
  wa: {
    1: { word: "挖", pinyin: "wā", example: "挖掘" },
    2: { word: "娃", pinyin: "wá", example: "娃娃" },
    3: { word: "瓦", pinyin: "wǎ", example: "瓦片" },
    4: { word: "袜", pinyin: "wà", example: "袜子" },
  },
  wo: {
    1: { word: "窝", pinyin: "wō", example: "鸟窝" },
    3: { word: "我", pinyin: "wǒ", example: "我们" },
    4: { word: "握", pinyin: "wò", example: "握手" },
  },
  wu: {
    1: { word: "屋", pinyin: "wū", example: "房屋" },
    2: { word: "无", pinyin: "wú", example: "没有" },
    3: { word: "五", pinyin: "wǔ", example: "数字五" },
    4: { word: "物", pinyin: "wù", example: "动物" },
  },
  wei: {
    1: { word: "威", pinyin: "wēi", example: "威风" },
    2: { word: "围", pinyin: "wéi", example: "周围" },
    3: { word: "尾", pinyin: "wěi", example: "尾巴" },
    4: { word: "位", pinyin: "wèi", example: "位置" },
  },
  wan: {
    1: { word: "弯", pinyin: "wān", example: "弯曲" },
    2: { word: "玩", pinyin: "wán", example: "玩耍" },
    3: { word: "碗", pinyin: "wǎn", example: "饭碗" },
    4: { word: "万", pinyin: "wàn", example: "千万" },
  },
  wen: {
    1: { word: "温", pinyin: "wēn", example: "温暖" },
    2: { word: "文", pinyin: "wén", example: "文化" },
    3: { word: "吻", pinyin: "wěn", example: "亲吻" },
    4: { word: "问", pinyin: "wèn", example: "问题" },
  },
  wang: {
    1: { word: "汪", pinyin: "wāng", example: "汪汪" },
    2: { word: "王", pinyin: "wáng", example: "国王" },
    3: { word: "网", pinyin: "wǎng", example: "上网" },
    4: { word: "忘", pinyin: "wàng", example: "忘记" },
  },
  weng: {
    1: { word: "翁", pinyin: "wēng", example: "老翁" },
    3: { word: "嗡", pinyin: "wěng", example: "嗡嗡" },
    4: { word: "瓮", pinyin: "wèng", example: "瓮" },
  },

  // 复韵母组合示例 (扩充)
  bei: {
    1: { word: "背", pinyin: "bēi", example: "背书包" },
    3: { word: "北", pinyin: "běi", example: "北京" },
    4: { word: "背", pinyin: "bèi", example: "后背" },
  },
  bao: {
    1: { word: "包", pinyin: "bāo", example: "书包" },
    3: { word: "宝", pinyin: "bǎo", example: "宝贝" },
    4: { word: "抱", pinyin: "bào", example: "拥抱" },
  },
  ban: {
    1: { word: "班", pinyin: "bān", example: "班级" },
    3: { word: "板", pinyin: "bǎn", example: "木板" },
    4: { word: "半", pinyin: "bàn", example: "一半" },
  },
  bang: {
    1: { word: "帮", pinyin: "bāng", example: "帮忙" },
    3: { word: "榜", pinyin: "bǎng", example: "榜样" },
    4: { word: "棒", pinyin: "bàng", example: "真棒" },
  },
  ben: {
    1: { word: "奔", pinyin: "bēn", example: "奔跑" },
    3: { word: "本", pinyin: "běn", example: "书本" },
    4: { word: "笨", pinyin: "bèn", example: "笨蛋" },
  },
  beng: {
    1: { word: "崩", pinyin: "bēng", example: "崩溃" },
    3: { word: "绷", pinyin: "běng", example: "绷带" },
    4: { word: "蹦", pinyin: "bèng", example: "蹦跳" },
  },
  bian: {
    1: { word: "边", pinyin: "biān", example: "旁边" },
    3: { word: "扁", pinyin: "biǎn", example: "扁担" },
    4: { word: "变", pinyin: "biàn", example: "变化" },
  },
  biao: {
    1: { word: "标", pinyin: "biāo", example: "标准" },
    3: { word: "表", pinyin: "biǎo", example: "手表" },
  },
  bie: {
    1: { word: "憋", pinyin: "biē", example: "憋气" },
    2: { word: "别", pinyin: "bié", example: "别人" },
  },
  bin: { 1: { word: "宾", pinyin: "bīn", example: "宾客" } },
  bing: {
    1: { word: "冰", pinyin: "bīng", example: "冰块" },
    3: { word: "饼", pinyin: "bǐng", example: "饼干" },
    4: { word: "病", pinyin: "bìng", example: "生病" },
  },

  pai: {
    1: { word: "拍", pinyin: "pāi", example: "拍手" },
    2: { word: "排", pinyin: "pái", example: "排队" },
    4: { word: "派", pinyin: "pài", example: "派对" },
  },
  pao: {
    1: { word: "抛", pinyin: "pāo", example: "抛弃" },
    2: { word: "袍", pinyin: "páo", example: "长袍" },
    3: { word: "跑", pinyin: "pǎo", example: "跑步" },
    4: { word: "炮", pinyin: "pào", example: "鞭炮" },
  },
  pei: {
    1: { word: "胚", pinyin: "pēi", example: "胚芽" },
    2: { word: "陪", pinyin: "péi", example: "陪伴" },
    4: { word: "配", pinyin: "pèi", example: "配合" },
  },
  pen: {
    1: { word: "喷", pinyin: "pēn", example: "喷泉" },
    2: { word: "盆", pinyin: "pén", example: "脸盆" },
  },
  peng: {
    1: { word: "烹", pinyin: "pēng", example: "烹饪" },
    2: { word: "朋", pinyin: "péng", example: "朋友" },
    3: { word: "捧", pinyin: "pěng", example: "捧场" },
    4: { word: "碰", pinyin: "pèng", example: "碰见" },
  },
  pi: {
    1: { word: "批", pinyin: "pī", example: "批评" },
    2: { word: "皮", pinyin: "pí", example: "皮肤" },
    3: { word: "匹", pinyin: "pǐ", example: "马匹" },
    4: { word: "屁", pinyin: "pì", example: "放屁" },
  },
  pian: {
    1: { word: "篇", pinyin: "piān", example: "篇章" },
    2: { word: "便", pinyin: "pián", example: "便宜" },
    3: { word: "片", pinyin: "piàn", example: "照片" },
  }, // note: pian4 is usually pian
  piao: {
    1: { word: "飘", pinyin: "piāo", example: "飘扬" },
    2: { word: "瓢", pinyin: "piáo", example: "水瓢" },
    3: { word: "漂", pinyin: "piǎo", example: "漂白" },
    4: { word: "票", pinyin: "piào", example: "车票" },
  },
  pin: {
    1: { word: "拼", pinyin: "pīn", example: "拼音" },
    2: { word: "贫", pinyin: "pín", example: "贫穷" },
    3: { word: "品", pinyin: "pǐn", example: "品尝" },
  },
  ping: {
    1: { word: "乒", pinyin: "pīng", example: "乒乓" },
    2: { word: "平", pinyin: "píng", example: "平安" },
  },

  mai: {
    2: { word: "埋", pinyin: "mái", example: "埋藏" },
    3: { word: "买", pinyin: "mǎi", example: "买卖" },
    4: { word: "卖", pinyin: "mài", example: "售卖" },
  },
  man: {
    1: { word: "蛮", pinyin: "mān", example: "蛮横" },
    2: { word: "满", pinyin: "mán", example: "充满" },
    3: { word: "满", pinyin: "mǎn", example: "满意" },
    4: { word: "慢", pinyin: "màn", example: "缓慢" },
  },
  mang: {
    2: { word: "忙", pinyin: "máng", example: "帮忙" },
    3: { word: "盲", pinyin: "mǎng", example: "文盲" },
  }, // mang3 is rare, mang2 is mang
  mao: {
    1: { word: "猫", pinyin: "māo", example: "小猫" },
    2: { word: "毛", pinyin: "máo", example: "毛巾" },
    3: { word: "卯", pinyin: "mǎo", example: "卯时" },
    4: { word: "帽", pinyin: "mào", example: "帽子" },
  },
  mei: {
    2: { word: "梅", pinyin: "méi", example: "梅花" },
    3: { word: "美", pinyin: "měi", example: "美丽" },
    4: { word: "妹", pinyin: "mèi", example: "妹妹" },
  },
  men: {
    2: { word: "门", pinyin: "mén", example: "大门" },
    4: { word: "闷", pinyin: "mèn", example: "闷热" },
  },
  meng: {
    1: { word: "蒙", pinyin: "mēng", example: "蒙骗" },
    2: { word: "萌", pinyin: "méng", example: "萌芽" },
    3: { word: "猛", pinyin: "měng", example: "凶猛" },
    4: { word: "梦", pinyin: "mèng", example: "做梦" },
  },
  mian: {
    2: { word: "棉", pinyin: "mián", example: "棉花" },
    3: { word: "免", pinyin: "miǎn", example: "免费" },
    4: { word: "面", pinyin: "miàn", example: "面条" },
  },
  miao: {
    1: { word: "喵", pinyin: "miāo", example: "喵喵" },
    2: { word: "苗", pinyin: "miáo", example: "禾苗" },
    3: { word: "秒", pinyin: "miǎo", example: "秒表" },
    4: { word: "妙", pinyin: "miào", example: "奇妙" },
  },
  mie: {
    1: { word: "咩", pinyin: "miē", example: "羊叫" },
    4: { word: "灭", pinyin: "miè", example: "消灭" },
  },
  min: {
    2: { word: "民", pinyin: "mín", example: "人民" },
    3: { word: "敏", pinyin: "mǐn", example: "敏感" },
  },
  ming: {
    2: { word: "明", pinyin: "míng", example: "明天" },
    3: { word: "酩", pinyin: "mǐng", example: "酩酊" },
    4: { word: "命", pinyin: "mìng", example: "命令" },
  },
  mou: {
    2: { word: "谋", pinyin: "móu", example: "计谋" },
    3: { word: "某", pinyin: "mǒu", example: "某人" },
  },

  fan: {
    1: { word: "翻", pinyin: "fān", example: "翻书" },
    2: { word: "凡", pinyin: "fán", example: "平凡" },
    3: { word: "反", pinyin: "fǎn", example: "反对" },
    4: { word: "饭", pinyin: "fàn", example: "吃饭" },
  },
  fang: {
    1: { word: "方", pinyin: "fāng", example: "方形" },
    2: { word: "房", pinyin: "fáng", example: "房子" },
    3: { word: "访", pinyin: "fǎng", example: "访问" },
    4: { word: "放", pinyin: "fàng", example: "放心" },
  },
  fei: {
    1: { word: "飞", pinyin: "fēi", example: "飞机" },
    2: { word: "肥", pinyin: "féi", example: "肥胖" },
    3: { word: "匪", pinyin: "fěi", example: "土匪" },
    4: { word: "费", pinyin: "fèi", example: "费用" },
  },
  fen: {
    1: { word: "分", pinyin: "fēn", example: "分数" },
    2: { word: "坟", pinyin: "fén", example: "坟墓" },
    3: { word: "粉", pinyin: "fěn", example: "粉色" },
    4: { word: "份", pinyin: "fèn", example: "一份" },
  },
  feng: {
    1: { word: "风", pinyin: "fēng", example: "大风" },
    2: { word: "缝", pinyin: "féng", example: "缝补" },
    3: { word: "讽", pinyin: "fěng", example: "讽刺" },
    4: { word: "凤", pinyin: "fèng", example: "凤凰" },
  },
  fo: { 2: { word: "佛", pinyin: "fó", example: "佛像" } },
  fou: { 3: { word: "否", pinyin: "fǒu", example: "否定" } },

  dai: {
    1: { word: "呆", pinyin: "dāi", example: "发呆" },
    3: { word: "逮", pinyin: "dǎi", example: "逮捕" },
    4: { word: "带", pinyin: "dài", example: "皮带" },
  },
  dan: {
    1: { word: "单", pinyin: "dān", example: "简单" },
    3: { word: "胆", pinyin: "dǎn", example: "胆小" },
    4: { word: "蛋", pinyin: "dàn", example: "鸡蛋" },
  },
  dang: {
    1: { word: "当", pinyin: "dāng", example: "当时" },
    3: { word: "挡", pinyin: "dǎng", example: "挡住" },
    4: { word: "荡", pinyin: "dàng", example: "飘荡" },
  },
  dao: {
    1: { word: "刀", pinyin: "dāo", example: "剪刀" },
    3: { word: "导", pinyin: "dǎo", example: "向导" },
    4: { word: "到", pinyin: "dào", example: "到达" },
  },
  deng: {
    1: { word: "灯", pinyin: "dēng", example: "电灯" },
    3: { word: "等", pinyin: "děng", example: "等待" },
    4: { word: "凳", pinyin: "dèng", example: "凳子" },
  },
  di: {
    1: { word: "低", pinyin: "dī", example: "高低" },
    2: { word: "敌", pinyin: "dí", example: "敌人" },
    3: { word: "底", pinyin: "dǐ", example: "底下" },
    4: { word: "弟", pinyin: "dì", example: "弟弟" },
  },
  dian: {
    1: { word: "颠", pinyin: "diān", example: "颠倒" },
    3: { word: "点", pinyin: "diǎn", example: "点心" },
    4: { word: "电", pinyin: "diàn", example: "电话" },
  },
  diao: {
    1: { word: "雕", pinyin: "diāo", example: "雕刻" },
    3: { word: "鸟", pinyin: "diǎo", example: "小鸟" },
    4: { word: "钓", pinyin: "diào", example: "钓鱼" },
  },
  die: {
    1: { word: "跌", pinyin: "diē", example: "跌倒" },
    2: { word: "叠", pinyin: "dié", example: "折叠" },
    3: { word: "爹", pinyin: "diě", example: "爹爹" },
  },
  ding: {
    1: { word: "钉", pinyin: "dīng", example: "钉子" },
    3: { word: "顶", pinyin: "dǐng", example: "头顶" },
    4: { word: "定", pinyin: "dìng", example: "决定" },
  },
  diu: { 1: { word: "丢", pinyin: "diū", example: "丢弃" } },
  dong: {
    1: { word: "东", pinyin: "dōng", example: "东方" },
    3: { word: "懂", pinyin: "dǒng", example: "懂事" },
    4: { word: "动", pinyin: "dòng", example: "运动" },
  },
  dou: {
    1: { word: "兜", pinyin: "dōu", example: "衣兜" },
    3: { word: "斗", pinyin: "dǒu", example: "北斗" },
    4: { word: "豆", pinyin: "dòu", example: "豆子" },
  },
  du: {
    1: { word: "嘟", pinyin: "dū", example: "嘟嘴" },
    2: { word: "毒", pinyin: "dú", example: "有毒" },
    3: { word: "肚", pinyin: "dǔ", example: "肚子" },
    4: { word: "度", pinyin: "dù", example: "温度" },
  },
  duan: {
    1: { word: "端", pinyin: "duān", example: "端正" },
    3: { word: "短", pinyin: "duǎn", example: "长短" },
    4: { word: "段", pinyin: "duàn", example: "片段" },
  },
  dui: {
    1: { word: "堆", pinyin: "duī", example: "土堆" },
    3: { word: "队", pinyin: "duì", example: "排队" },
    4: { word: "对", pinyin: "duì", example: "不对" },
  },
  dun: {
    1: { word: "蹲", pinyin: "dūn", example: "蹲下" },
    3: { word: "盹", pinyin: "dǔn", example: "打盹" },
    4: { word: "盾", pinyin: "dùn", example: "盾牌" },
  },
  duo: {
    1: { word: "多", pinyin: "duō", example: "多少" },
    2: { word: "夺", pinyin: "duó", example: "争夺" },
    3: { word: "朵", pinyin: "duǒ", example: "花朵" },
    4: { word: "跺", pinyin: "duò", example: "跺脚" },
  },

  tai: {
    1: { word: "胎", pinyin: "tāi", example: "轮胎" },
    2: { word: "台", pinyin: "tái", example: "讲台" },
    4: { word: "太", pinyin: "tài", example: "太阳" },
  },
  tan: {
    1: { word: "贪", pinyin: "tān", example: "贪心" },
    2: { word: "谈", pinyin: "tán", example: "谈话" },
    3: { word: "毯", pinyin: "tǎn", example: "地毯" },
    4: { word: "探", pinyin: "tàn", example: "探险" },
  },
  tang: {
    1: { word: "汤", pinyin: "tāng", example: "喝汤" },
    2: { word: "糖", pinyin: "táng", example: "糖果" },
    3: { word: "躺", pinyin: "tǎng", example: "躺下" },
    4: { word: "烫", pinyin: "tàng", example: "滚烫" },
  },
  tao: {
    1: { word: "涛", pinyin: "tāo", example: "波涛" },
    2: { word: "桃", pinyin: "táo", example: "桃子" },
    3: { word: "讨", pinyin: "tǎo", example: "讨论" },
    4: { word: "套", pinyin: "tào", example: "手套" },
  },
  te: { 4: { word: "特", pinyin: "tè", example: "特别" } },
  teng: { 2: { word: "疼", pinyin: "téng", example: "疼痛" } },
  ti: {
    1: { word: "梯", pinyin: "tī", example: "楼梯" },
    2: { word: "提", pinyin: "tí", example: "提问" },
    3: { word: "体", pinyin: "tǐ", example: "身体" },
    4: { word: "替", pinyin: "tì", example: "代替" },
  },
  tian: {
    1: { word: "天", pinyin: "tiān", example: "天空" },
    2: { word: "田", pinyin: "tián", example: "田野" },
    3: { word: "舔", pinyin: "tiǎn", example: "舔食" },
  },
  tiao: {
    1: { word: "挑", pinyin: "tiāo", example: "挑选" },
    2: { word: "条", pinyin: "tiáo", example: "面条" },
    3: { word: "挑", pinyin: "tiǎo", example: "挑战" },
    4: { word: "跳", pinyin: "tiào", example: "跳高" },
  },
  tie: {
    1: { word: "贴", pinyin: "tiē", example: "粘贴" },
    3: { word: "铁", pinyin: "tiě", example: "铁路" },
  },
  ting: {
    1: { word: "听", pinyin: "tīng", example: "听说" },
    2: { word: "停", pinyin: "tíng", example: "停止" },
    3: { word: "挺", pinyin: "tǐng", example: "挺好" },
  },
  tong: {
    1: { word: "通", pinyin: "tōng", example: "通过" },
    2: { word: "同", pinyin: "tóng", example: "同学" },
    3: { word: "桶", pinyin: "tǒng", example: "水桶" },
    4: { word: "痛", pinyin: "tòng", example: "痛苦" },
  },
  tou: {
    1: { word: "偷", pinyin: "tōu", example: "小偷" },
    2: { word: "头", pinyin: "tóu", example: "头发" },
    4: { word: "透", pinyin: "tòu", example: "透明" },
  },
  tuan: { 2: { word: "团", pinyin: "tuán", example: "团圆" } },
  tui: {
    1: { word: "推", pinyin: "tuī", example: "推开" },
    2: { word: "颓", pinyin: "tuí", example: "颓废" },
    3: { word: "腿", pinyin: "tuǐ", example: "大腿" },
    4: { word: "退", pinyin: "tuì", example: "后退" },
  },
  tun: {
    1: { word: "吞", pinyin: "tūn", example: "吞咽" },
    2: { word: "屯", pinyin: "tún", example: "屯兵" },
  },
  tuo: {
    1: { word: "托", pinyin: "tuō", example: "托运" },
    2: { word: "驼", pinyin: "tuó", example: "骆驼" },
    3: { word: "妥", pinyin: "tuǒ", example: "妥当" },
  },

  nai: {
    3: { word: "奶", pinyin: "nǎi", example: "牛奶" },
    4: { word: "耐", pinyin: "nài", example: "耐心" },
  },
  nan: {
    2: { word: "南", pinyin: "nán", example: "南方" },
    3: { word: "难", pinyin: "nǎn", example: "困难" },
    4: { word: "难", pinyin: "nàn", example: "灾难" },
  },
  nang: { 2: { word: "囊", pinyin: "náng", example: "锦囊" } },
  nao: {
    2: { word: "挠", pinyin: "náo", example: "挠痒" },
    3: { word: "脑", pinyin: "nǎo", example: "电脑" },
    4: { word: "闹", pinyin: "nào", example: "热闹" },
  },
  ne: {
    0: { word: "呢", pinyin: "ne", example: "你呢" },
    4: { word: "呐", pinyin: "nè", example: "呐喊" },
  }, // 0声 usually not in tones array (1-4), map to closest or skip
  nei: {
    3: { word: "馁", pinyin: "něi", example: "气馁" },
    4: { word: "内", pinyin: "nèi", example: "内部" },
  },
  nen: { 4: { word: "嫩", pinyin: "nèn", example: "鲜嫩" } },
  neng: { 2: { word: "能", pinyin: "néng", example: "能力" } },
  nian: {
    2: { word: "年", pinyin: "nián", example: "过年" },
    3: { word: "捻", pinyin: "niǎn", example: "捻线" },
    4: { word: "念", pinyin: "niàn", example: "想念" },
  },
  niang: {
    2: { word: "娘", pinyin: "niáng", example: "姑娘" },
    4: { word: "酿", pinyin: "niàng", example: "酿酒" },
  },
  niao: {
    3: { word: "鸟", pinyin: "niǎo", example: "小鸟" },
    4: { word: "尿", pinyin: "niào", example: "尿布" },
  },
  nie: {
    1: { word: "捏", pinyin: "niē", example: "捏泥人" },
    4: { word: "聂", pinyin: "niè", example: "姓聂" },
  },
  nin: { 2: { word: "您", pinyin: "nín", example: "您好" } },
  ning: {
    2: { word: "宁", pinyin: "níng", example: "安宁" },
    3: { word: "拧", pinyin: "nǐng", example: "拧干" },
    4: { word: "泞", pinyin: "nìng", example: "泥泞" },
  },
  niu: {
    1: { word: "妞", pinyin: "niū", example: "小妞" },
    2: { word: "牛", pinyin: "niú", example: "黄牛" },
    3: { word: "扭", pinyin: "niǔ", example: "扭动" },
  },
  nong: {
    2: { word: "农", pinyin: "nóng", example: "农民" },
    4: { word: "弄", pinyin: "nòng", example: "摆弄" },
  },
  nou: { 4: { word: "耨", pinyin: "nòu", example: "耨" } },
  nu: {
    2: { word: "奴", pinyin: "nú", example: "奴隶" },
    3: { word: "努", pinyin: "nǔ", example: "努力" },
    4: { word: "怒", pinyin: "nù", example: "愤怒" },
  },
  nuan: { 3: { word: "暖", pinyin: "nuǎn", example: "温暖" } },
  nue: { 4: { word: "虐", pinyin: "nüè", example: "虐待" } },
  nuo: {
    2: { word: "挪", pinyin: "nuó", example: "挪动" },
    4: { word: "诺", pinyin: "nuò", example: "承诺" },
  },

  lai: {
    2: { word: "来", pinyin: "lái", example: "来到" },
    4: { word: "赖", pinyin: "lài", example: "依赖" },
  },
  lan: {
    2: { word: "蓝", pinyin: "lán", example: "蓝色" },
    3: { word: "懒", pinyin: "lǎn", example: "懒惰" },
    4: { word: "烂", pinyin: "làn", example: "灿烂" },
  },
  lang: {
    2: { word: "狼", pinyin: "láng", example: "野狼" },
    3: { word: "朗", pinyin: "lǎng", example: "朗读" },
    4: { word: "浪", pinyin: "làng", example: "海浪" },
  },
  lao: {
    1: { word: "捞", pinyin: "lāo", example: "打捞" },
    2: { word: "劳", pinyin: "láo", example: "劳动" },
    3: { word: "老", pinyin: "lǎo", example: "老人" },
    4: { word: "烙", pinyin: "lào", example: "烙饼" },
  },
  lei: {
    2: { word: "雷", pinyin: "léi", example: "雷雨" },
    3: { word: "累", pinyin: "lěi", example: "积累" },
    4: { word: "类", pinyin: "lèi", example: "种类" },
  },
  leng: { 3: { word: "冷", pinyin: "lěng", example: "寒冷" } },
  lia: { 3: { word: "俩", pinyin: "liǎ", example: "哥俩" } },
  lian: {
    2: { word: "连", pinyin: "lián", example: "连接" },
    3: { word: "脸", pinyin: "liǎn", example: "笑脸" },
    4: { word: "练", pinyin: "liàn", example: "练习" },
  },
  liang: {
    2: { word: "良", pinyin: "liáng", example: "良好" },
    3: { word: "两", pinyin: "liǎng", example: "两个" },
    4: { word: "亮", pinyin: "liàng", example: "明亮" },
  },
  liao: {
    2: { word: "聊", pinyin: "liáo", example: "聊天" },
    3: { word: "了", pinyin: "liǎo", example: "了解" },
    4: { word: "料", pinyin: "liào", example: "材料" },
  },
  lie: {
    1: { word: "咧", pinyin: "liē", example: "咧嘴" },
    3: { word: "咧", pinyin: "liě", example: "咧" },
    4: { word: "列", pinyin: "liè", example: "列车" },
  },
  lin: {
    1: { word: "林", pinyin: "lín", example: "树林" },
    2: { word: "临", pinyin: "lín", example: "面临" },
    3: { word: "凛", pinyin: "lǐn", example: "凛冽" },
    4: { word: "吝", pinyin: "lìn", example: "吝啬" },
  },
  ling: {
    2: { word: "零", pinyin: "líng", example: "零食" },
    3: { word: "领", pinyin: "lǐng", example: "领带" },
    4: { word: "令", pinyin: "lìng", example: "命令" },
  },
  liu: {
    1: { word: "溜", pinyin: "liū", example: "溜冰" },
    2: { word: "流", pinyin: "liú", example: "流动" },
    3: { word: "柳", pinyin: "liǔ", example: "柳树" },
    4: { word: "六", pinyin: "liù", example: "数字六" },
  },
  long: {
    2: { word: "龙", pinyin: "lóng", example: "恐龙" },
    3: { word: "笼", pinyin: "lǒng", example: "灯笼" },
  },
  lou: {
    1: { word: "搂", pinyin: "lōu", example: "搂抱" },
    2: { word: "楼", pinyin: "lóu", example: "大楼" },
    3: { word: "搂", pinyin: "lǒu", example: "搂" },
    4: { word: "漏", pinyin: "lòu", example: "漏水" },
  },
  lu: {
    1: { word: "撸", pinyin: "lū", example: "撸袖子" },
    2: { word: "卢", pinyin: "lú", example: "葫芦" },
    3: { word: "鲁", pinyin: "lǔ", example: "鲁莽" },
    4: { word: "路", pinyin: "lù", example: "马路" },
  },
  luan: {
    2: { word: "峦", pinyin: "luán", example: "山峦" },
    3: { word: "卵", pinyin: "luǎn", example: "产卵" },
    4: { word: "乱", pinyin: "luàn", example: "乱画" },
  },
  lun: {
    1: { word: "抡", pinyin: "lūn", example: "抡锤" },
    2: { word: "轮", pinyin: "lún", example: "车轮" },
    4: { word: "论", pinyin: "lùn", example: "讨论" },
  },
  luo: {
    1: { word: "啰", pinyin: "luō", example: "啰嗦" },
    2: { word: "罗", pinyin: "luó", example: "罗列" },
    3: { word: "裸", pinyin: "luǒ", example: "裸露" },
    4: { word: "落", pinyin: "luò", example: "落叶" },
  },

  hei: { 1: { word: "黑", pinyin: "hēi", example: "黑色" } },
  gui: {
    1: { word: "龟", pinyin: "guī", example: "乌龟" },
    3: { word: "鬼", pinyin: "guǐ", example: "鬼怪" },
    4: { word: "贵", pinyin: "guì", example: "宝贵" },
  },
  jie: {
    1: { word: "接", pinyin: "jiē", example: "接受" },
    2: { word: "结", pinyin: "jié", example: "结实" },
    3: { word: "姐", pinyin: "jiě", example: "姐姐" },
    4: { word: "借", pinyin: "jiè", example: "借书" },
  },
  que: {
    1: { word: "缺", pinyin: "quē", example: "缺点" },
    2: { word: "瘸", pinyin: "qué", example: "瘸腿" },
    4: { word: "确", pinyin: "què", example: "确定" },
  },
  xue: {
    1: { word: "薛", pinyin: "xuē", example: "姓薛" },
    2: { word: "学", pinyin: "xué", example: "学习" },
    3: { word: "雪", pinyin: "xuě", example: "下雪" },
  },
  er: {
    2: { word: "儿", pinyin: "ér", example: "儿童" },
    3: { word: "耳", pinyin: "ěr", example: "耳朵" },
    4: { word: "二", pinyin: "èr", example: "数字二" },
  },
  cun: {
    1: { word: "村", pinyin: "cūn", example: "农村" },
    2: { word: "存", pinyin: "cún", example: "保存" },
    4: { word: "寸", pinyin: "cùn", example: "尺寸" },
  },
  qun: {
    2: { word: "群", pinyin: "qún", example: "群众" },
    3: { word: "裙", pinyin: "qún", example: "裙子" },
  }, // 裙子通常2声
};
