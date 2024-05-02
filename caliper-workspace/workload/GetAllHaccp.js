'use strict';

const { WorkloadModuleBase } = require('@hyperledger/caliper-core');
class MyWorkload extends WorkloadModuleBase {
    constructor() {
        super();
    }   
    
    async submitTransaction() {
        const myArgs = {
            contractId: 'basic',
            contractFunction: 'GetAllHaccp',
            invokerIdentity: 'User1',
            readOnly: true
        };
        await this.sutAdapter.sendRequests(myArgs);
    }
}    
function createWorkloadModule() {
    return new MyWorkload();
}
module.exports.createWorkloadModule = createWorkloadModule;


