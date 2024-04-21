export namespace api {
	
	export class EquQueryCond {
	    rarity: number[];
	    types: string[];
	    part_set_type: number;
	
	    static createFrom(source: any = {}) {
	        return new EquQueryCond(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.rarity = source["rarity"];
	        this.types = source["types"];
	        this.part_set_type = source["part_set_type"];
	    }
	}
	export class SkillDescription {
	    tp: number;
	    idx: number;
	    rate: number;
	
	    static createFrom(source: any = {}) {
	        return new SkillDescription(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.tp = source["tp"];
	        this.idx = source["idx"];
	        this.rate = source["rate"];
	    }
	}
	export class SkillData {
	    name: string;
	    code: number;
	    job_str: string;
	    static_data: number[];
	    level_data: number[][];
	    level_data_num: number;
	    description: SkillDescription[];
	    description_template: string;
	
	    static createFrom(source: any = {}) {
	        return new SkillData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.code = source["code"];
	        this.job_str = source["job_str"];
	        this.static_data = source["static_data"];
	        this.level_data = source["level_data"];
	        this.level_data_num = source["level_data_num"];
	        this.description = this.convertValues(source["description"], SkillDescription);
	        this.description_template = source["description_template"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace data_loader {
	
	export class CheckUpdateResult {
	    key: string;
	    name: string;
	    need_update: boolean;
	
	    static createFrom(source: any = {}) {
	        return new CheckUpdateResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.key = source["key"];
	        this.name = source["name"];
	        this.need_update = source["need_update"];
	    }
	}

}

export namespace model {
	
	export class Equipment {
	    code: number;
	    name: string;
	    path: string;
	    rarity: number;
	    type: string;
	    part_set: number;
	    mini_level: number;
	
	    static createFrom(source: any = {}) {
	        return new Equipment(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.code = source["code"];
	        this.name = source["name"];
	        this.path = source["path"];
	        this.rarity = source["rarity"];
	        this.type = source["type"];
	        this.part_set = source["part_set"];
	        this.mini_level = source["mini_level"];
	    }
	}
	export class Skill {
	    id: number;
	    job: number;
	    code: number;
	    name: string;
	    path: string;
	
	    static createFrom(source: any = {}) {
	        return new Skill(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.job = source["job"];
	        this.code = source["code"];
	        this.name = source["name"];
	        this.path = source["path"];
	    }
	}

}

export namespace script {
	
	export class Action {
	    type: string;
	    args: string[];
	
	    static createFrom(source: any = {}) {
	        return new Action(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.type = source["type"];
	        this.args = source["args"];
	    }
	}

}

export namespace setting {
	
	export class SettingStruct {
	    mode: string;
	    target: string;
	    concurrent_count: number;
	
	    static createFrom(source: any = {}) {
	        return new SettingStruct(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.mode = source["mode"];
	        this.target = source["target"];
	        this.concurrent_count = source["concurrent_count"];
	    }
	}

}

export namespace world_drop {
	
	export class WorldDropItem {
	    id: number;
	    level: number;
	    rate: number;
	
	    static createFrom(source: any = {}) {
	        return new WorldDropItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.level = source["level"];
	        this.rate = source["rate"];
	    }
	}

}

