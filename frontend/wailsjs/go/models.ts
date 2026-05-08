export namespace main {
	
	export class Article {
	    id: number;
	    title: string;
	    content: string;
	    created_at: string;
	
	    static createFrom(source: any = {}) {
	        return new Article(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.title = source["title"];
	        this.content = source["content"];
	        this.created_at = source["created_at"];
	    }
	}
	export class TranslateResult {
	    word: string;
	    translation: string;
	    phonetic: string;
	    cached: boolean;
	
	    static createFrom(source: any = {}) {
	        return new TranslateResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.word = source["word"];
	        this.translation = source["translation"];
	        this.phonetic = source["phonetic"];
	        this.cached = source["cached"];
	    }
	}
	export class WordBookItem {
	    id: number;
	    word: string;
	    translation: string;
	    phonetic: string;
	    note: string;
	    reviewed: number;
	    last_review: string;
	    created_at: string;
	
	    static createFrom(source: any = {}) {
	        return new WordBookItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.word = source["word"];
	        this.translation = source["translation"];
	        this.phonetic = source["phonetic"];
	        this.note = source["note"];
	        this.reviewed = source["reviewed"];
	        this.last_review = source["last_review"];
	        this.created_at = source["created_at"];
	    }
	}

}

