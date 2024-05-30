import {model} from "../../wailsjs/go/models";
import {GetJobInfo, GetJobSkill} from "../../wailsjs/go/api/App";

import Skill = model.Skill;
import Job = model.Job;

class Storage {
  private readonly _skillMap: Map<number, Map<number, Skill>> = new Map()
  private readonly _jobMap: Map<number, Array<Job>> = new Map()
  private readonly _jobStr2Code: Map<string, number> = new Map()
  private readonly loading = {
    job: false,
    skill: false,
  }

  async loadSkillMap() {
    if (this._skillMap.size > 0 || this.loading.skill) {
      return
    }
    this.loading.skill = true
    try {
      const skills = await GetJobSkill(0)
      skills.forEach((skl) => {
        if (!this._skillMap.has(skl.job)) {
          this._skillMap.set(skl.job, new Map())
        }
        this._skillMap.get(skl.job)!.set(skl.code, skl)
      })
    } finally {
      this.loading.skill = false
    }
  }

  async loadJobMap() {
    if (this._jobMap.size > 0 || this.loading.job) {
      return
    }
    this.loading.job = true
    try {
      const jobs = await GetJobInfo()
      jobs.forEach((item) => {
        if (!this._jobMap.has(item.code)) {
          this._jobMap.set(item.code, [item])
        } else {
          this._jobMap.get(item.code)!.push(item)
        }
        this._jobStr2Code.set(item.job, item.code)
      })
    } finally {
      this.loading.job = false
    }
  }

  getSkill(jobId: number, skillId: number) {
    return this._skillMap.get(jobId)?.get(skillId)
  }

  getJobMap() {
    return this._jobMap
  }

  getJobName(jobId: number) {
    const jobs = this._jobMap.get(jobId)
    return jobs ? jobs[0].job_name : String(jobId)
  }

  getJobId(job: string) {
    console.log(this._jobStr2Code, job)
    return this._jobStr2Code.get(job)
  }

  getJobStr(jobId: number) {
    for (const [k, v] of this._jobStr2Code) {
      if (jobId === v) {
        return k
      }
    }
    return 'all'
  }

  getJobNameByStr(job: string) {
    const jobId = this._jobStr2Code.get(job)
    if (jobId === undefined) {
      return job
    }
    return this.getJobName(jobId)
  }
}

const storage = new Storage()

export default storage

