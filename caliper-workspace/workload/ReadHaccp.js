'use strict';

const { WorkloadModuleBase } = require('@hyperledger/caliper-core');
class MyWorkload extends WorkloadModuleBase {
    constructor() {
        super();
        this.txIndex = 0;
        this.limitIndex =0;
    }   
    async initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext) {
        await super.initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext);
        this,this.limitIndex = this.roundArguments.assets;
    }
    async submitTransaction() {
        this.txIndex++;
        let id = 'fa' + this.txIndex 
        const myArgs = {
            contractId: 'basic',
            contractFunction: 'ReadHaccp',
            invokerIdentity: 'User1',
            contractArguments: [`Fa120230101`],
            readOnly: true
        };
        if (this.txIndex === this.limitIndex) {
            this.txIndex = 0;
        }
        await this.sutAdapter.sendRequests(myArgs);
    }
}    
function createWorkloadModule() {
    return new MyWorkload();
}
module.exports.createWorkloadModule = createWorkloadModule;