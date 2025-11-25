/// <reference types="vite/client" />
/// <reference types="vue" />
/// <reference types="vue-router" />
declare module '*.vue' { import { DefineComponent } from 'vue'; const component: DefineComponent<{}, {}, any>; export default component }
declare module '*.png' { const src: string; export default src }
declare module '*.jpg' { const src: string; export default src }
declare module '*.jpeg' { const src: string; export default src }
declare module '*.webp' { const src: string; export default src }
declare module '*.svg' { const src: string; export default src }
declare module '*.gif' { const src: string; export default src }
