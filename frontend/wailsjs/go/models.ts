export namespace main {
	
	export class SiteStatus {
	    url: string;
	    status: string;
	    message: string;
	    isOnline: boolean;
	
	    static createFrom(source: any = {}) {
	        return new SiteStatus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.url = source["url"];
	        this.status = source["status"];
	        this.message = source["message"];
	        this.isOnline = source["isOnline"];
	    }
	}

}

