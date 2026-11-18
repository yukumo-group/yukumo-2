declare module 'yukumo.js' {
  export type Voice =
    | 'dvd'
    | 'f1'
    | 'f2'
    | 'imd1'
    | 'jgr'
    | 'm1'
    | 'm2'
    | 'r1';

  export type AquesTalk2Voice =
    | 'f4'
    | 'rm3'
    | 'm5'
    | 'mf2'
    | 'mf1'
    | 'rm'
    | 'f1c'
    | 'm4b'
    | 'robo'
    | 'rb3'
    | 'rb2'
    | 'f3a'
    | 'huskey'
    | 'yukkuri'
    | 'defo1'
    | 'teto1'
    | 'momo1';

  export type AquesTalk10Voice = 'f1' | 'f2' | 'f3' | 'm1' | 'm2' | 'r1' | 'r2';

  export interface AqtkVoice {
    bas: number;
    spd: number;
    vol: number;
    pit: number;
    acc: number;
    lmd: number;
    fsc: number;
  }

  export interface Options {
    memorySize?: number;
  }

  export class AquesTalk1 {
    SetDevKey(key: string): number;
    run(koe: string, speed?: number): Uint8Array;
    destroy(): Promise<void>;
  }

  export class AquesTalk2 {
    run(koe: string, speed?: number): Uint8Array;
    destroy(): Promise<void>;
  }

  export class AquesTalk10 {
    SetDevKey(key: string): number;
    run(koe: string, speed?: number, voice?: AquesTalk10Voice | AqtkVoice): Uint8Array;
    destroy(): Promise<void>;
  }

  export function load(voice: Voice, options?: Options): Promise<AquesTalk1>;
  export function loadAquesTalk2(
    voice: AquesTalk2Voice,
    options?: Options,
  ): Promise<AquesTalk2>;
  export function loadAquesTalk10(
    voice: AquesTalk10Voice,
    options?: Options,
  ): Promise<AquesTalk10>;
}

declare module 'yukumo.js/lang/kanji2koe' {
  export interface Options {
    memorySize?: number;
  }

  export class AqKanji2Koe {
    setDevKey(key: string): number;
    convert(text: string): string;
    destroy(): Promise<void>;
  }

  export function loadAqKanji2Koe(options?: Options): Promise<AqKanji2Koe>;
}
