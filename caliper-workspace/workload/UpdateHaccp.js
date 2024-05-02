'use strict';
const crypto = require('crypto');
const { WorkloadModuleBase } = require('@hyperledger/caliper-core');
class MyWorkload extends WorkloadModuleBase {
    constructor() {
        super();
        this.txIndex = 0;   
    }   

    // async submitTransaction() {
    //     this.txIndex ++ ;
    //     let fa = 'fa' + this.txIndex 
    //     let hash = crypto.createHash('sha256')
    //     hash.update(fa);
    //     let mkroot = hash.digest('hex')
    //     let now = new Date();
    //     let Time = now.toISOString().replace(/T/, ' ').replace(/\..+/, '');
    //     const myArgs = {
    //         contractId: 'basic',
    //         contractFunction: 'UpdateHaccp',
    //         invokerIdentity: 'User1',
    //         contractArguments: [fa,mkroot,Time],
    //         readOnly: false
    //     };
    //     await this.sutAdapter.sendRequests(myArgs);
    // }
    async submitTransaction() {
        this.txIndex ++ ;
        let fa = 'fa' + this.txIndex 
        let hash = crypto.createHash('sha256')
        hash.update(fa);
        let mkroot = hash.digest('hex')
        const myArgs = {
            contractId: 'basic',
            contractFunction: 'UpdateHaccp',
            invokerIdentity: 'User1',
            contractArguments: [fa,mkroot],
            readOnly: false,
            timeout: 60
        };
        await this.sutAdapter.sendRequests(myArgs);
    }
    
    


}    
function createWorkloadModule() {
    return new MyWorkload();
}
module.exports.createWorkloadModule = createWorkloadModule;
