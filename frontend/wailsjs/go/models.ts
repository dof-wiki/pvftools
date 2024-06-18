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
	export class SearchResult {
	    id: number;
	    name: string;
	
	    static createFrom(source: any = {}) {
	        return new SearchResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
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

export namespace box {
	
	export class BoxItem {
	    id: number;
	    count: number;
	    rate: number;
	
	    static createFrom(source: any = {}) {
	        return new BoxItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.count = source["count"];
	        this.rate = source["rate"];
	    }
	}
	export class BoxItemGroup {
	    items: BoxItem[];
	
	    static createFrom(source: any = {}) {
	        return new BoxItemGroup(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.items = this.convertValues(source["items"], BoxItem);
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
	export class BoxParams {
	    box_type: number;
	    name: string;
	    explain: string;
	    grade: number;
	    rarity: number;
	    items: BoxItemGroup[];
	    tp: string;
	    category_names: string[];
	
	    static createFrom(source: any = {}) {
	        return new BoxParams(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.box_type = source["box_type"];
	        this.name = source["name"];
	        this.explain = source["explain"];
	        this.grade = source["grade"];
	        this.rarity = source["rarity"];
	        this.items = this.convertValues(source["items"], BoxItemGroup);
	        this.tp = source["tp"];
	        this.category_names = source["category_names"];
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
	
	export class CustomAttrTmpl {
	    id: number;
	    name: string;
	    tmpl: string;
	    desc: string;
	    choices: string;
	
	    static createFrom(source: any = {}) {
	        return new CustomAttrTmpl(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.tmpl = source["tmpl"];
	        this.desc = source["desc"];
	        this.choices = source["choices"];
	    }
	}
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
	export class Job {
	    id: number;
	    code: number;
	    job: string;
	    grow_type: number;
	    grow_type_name: string;
	    job_name: string;
	
	    static createFrom(source: any = {}) {
	        return new Job(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.code = source["code"];
	        this.job = source["job"];
	        this.grow_type = source["grow_type"];
	        this.grow_type_name = source["grow_type_name"];
	        this.job_name = source["job_name"];
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

export namespace proto {
	
	export class SkillIncrement {
	    job: string;
	    skill_id: number;
	    type: string;
	    data_type: string;
	    data_index: number;
	    increment_type: string;
	    value: number;
	
	    static createFrom(source: any = {}) {
	        return new SkillIncrement(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.job = source["job"];
	        this.skill_id = source["skill_id"];
	        this.type = source["type"];
	        this.data_type = source["data_type"];
	        this.data_index = source["data_index"];
	        this.increment_type = source["increment_type"];
	        this.value = source["value"];
	    }
	}
	export class BreathCheckSkillItem {
	    index: number;
	    increment_list: SkillIncrement[];
	    explain: string;
	
	    static createFrom(source: any = {}) {
	        return new BreathCheckSkillItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.index = source["index"];
	        this.increment_list = this.convertValues(source["increment_list"], SkillIncrement);
	        this.explain = source["explain"];
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
	export class BreathCheckItem {
	    job: number;
	    sub_job: number;
	    equ_type: string;
	    skill_list: BreathCheckSkillItem[];
	
	    static createFrom(source: any = {}) {
	        return new BreathCheckItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.job = source["job"];
	        this.sub_job = source["sub_job"];
	        this.equ_type = source["equ_type"];
	        this.skill_list = this.convertValues(source["skill_list"], BreathCheckSkillItem);
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
	export class BreathData {
	    id: number;
	    name: string;
	    probability1: number;
	    probability2: number;
	    check_list: BreathCheckItem[];
	    has_change: boolean;
	
	    static createFrom(source: any = {}) {
	        return new BreathData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.probability1 = source["probability1"];
	        this.probability2 = source["probability2"];
	        this.check_list = this.convertValues(source["check_list"], BreathCheckItem);
	        this.has_change = source["has_change"];
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
	export class UpgradeItem {
	    level: number;
	    fail_rate: number;
	    fail_op: number;
	    fail_op_value: number;
	    cost_id: number;
	    cost_count: number;
	
	    static createFrom(source: any = {}) {
	        return new UpgradeItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.level = source["level"];
	        this.fail_rate = source["fail_rate"];
	        this.fail_op = source["fail_op"];
	        this.fail_op_value = source["fail_op_value"];
	        this.cost_id = source["cost_id"];
	        this.cost_count = source["cost_count"];
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

