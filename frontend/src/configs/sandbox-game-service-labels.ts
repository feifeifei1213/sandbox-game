const vipServiceOperatingLabels = {
  taxPolicyNote: '一般企业 25%，支线机场贵宾 15%',
  materialGroupTitle: '2. 支付前期花费费用',
  productionAdjustmentGroupTitle: '3. 贵宾厅调整',
  humanResourceGroupTitle: '5. 人力资源管理',
  humanResourceLabel: '支付招聘、辞退、培训、待岗费用',
  salaryGroupTitle: '6. 贵宾厅空位开始服务，发放服务人员工资',
  salaryLabel: '支付贵宾厅的服务人员工资',
  researchAndManagementGroupTitle: '7. 新服务研创 / 8. 管理体系投入',
  newOrderReminder: '9. 下新服务订单',
  receivableGroupTitle: '10. 更新应收账款',
  receivableLabel: '如应收账款进入现金区做记录',
  receivableReminder: '所有应收向左移动一格',
  deliveryGroupTitle: '11. 服务验收交付',
  deliveryRevenueLabel: '服务订单销售额',
  deliveryCostLabel: '服务订单成本',
  managementStaffGroupTitle: '12. 支付管理人员费用',
  lineMaintenanceGroupTitle: '2. 支付贵宾厅年度维护费',
  lineMaintenanceNote: '1M / 间',
  assetGroupTitle: '3. 贵宾区资产',
  assetValueNote: '贵宾A / B / C区 价值（40M、32M、16M）',
  rentGroupTitle: '4. 贵宾区租金',
  rentValueNote: '贵宾A / B / C区 租金（4M、3M、2M）',
  residualGroupTitle: '5. 贵宾厅残值',
  depreciationGroupTitle: '6. 贵宾厅折旧',
  workInConstructionLabel: '未完工贵宾厅价值',
  marketProductFields: [
    { key: 'agencyInspectionTotal', label: '代办过检服务总价' },
    { key: 'twoCabinVipTotal', label: '两舱贵宾服务总价' },
    { key: 'businessVipTotal', label: '商务贵宾服务总价' },
    { key: 'memberCustomTotal', label: '会员定制服务总价' },
  ],
  materialFields: [
    { key: 'basicProduct', label: '代办过检服务' },
    { key: 'standardProduct', label: '两舱贵宾服务' },
    { key: 'precisionProduct', label: '商务贵宾服务' },
    { key: 'intelligentProduct', label: '会员定制服务' },
  ],
  productionLineRows: [
    { key: 'changeProduct', label: '变更服务', labelColspan: 2 },
    { key: 'dismantleCost', label: '贵宾厅拆除 1M / 间', labelColspan: 2 },
    { key: 'lineSale', label: '贵宾厅出售', labelColspan: 2 },
    { key: 'newLineInstall', label: '新贵宾厅建造', labelColspan: 2 },
    { key: 'constructionToFixed', label: '贵宾厅残值', groupLabel: '建成贵宾厅转固', groupRowspan: 2 },
    { key: 'newDepreciableAsset', label: '待折资产', useExistingGroupCell: true },
  ],
  researchFields: [
    { key: 'technologyResearch', label: '新服务研创' },
    { key: 'managementSystem', label: '质量环境健康费用' },
  ],
  progressUpdateReminder: '服务进度更新',
} as const

const productionOperatingLabels = {
  taxPolicyNote: '一般企业 25%，高新企业 15%',
  materialGroupTitle: '2. 支付材料费用',
  productionAdjustmentGroupTitle: '3. 生产线调整',
  humanResourceGroupTitle: '4. 人力资源',
  humanResourceLabel: '人力资源费用',
  salaryGroupTitle: '5. 工资与生产',
  salaryLabel: '工资与生产费用',
  researchAndManagementGroupTitle: '6. 研发与管理',
  newOrderReminder: '下新供应链订单',
  receivableGroupTitle: '7. 应收更新',
  receivableLabel: '应收回款',
  receivableReminder: '请将应收账款向左移动一格',
  deliveryGroupTitle: '8. 产品验收交付',
  deliveryRevenueLabel: '销售收入',
  deliveryCostLabel: '直接成本',
  managementStaffGroupTitle: '9. 支付管理人员费用',
  lineMaintenanceGroupTitle: '2. 支付生产线年度维护费',
  lineMaintenanceNote: '1M / 条',
  assetGroupTitle: '3. 生产厂房',
  assetValueNote: 'A / B / C 厂房价值（40M、32M、16M）',
  rentGroupTitle: '4. 生产厂房租金',
  rentValueNote: 'A / B / C 厂房租金（4M、3M、2M）',
  residualGroupTitle: '5. 生产线残值',
  depreciationGroupTitle: '6. 生产线折旧',
  workInConstructionLabel: '未完工生产线价值',
  marketProductFields: [
    { key: 'agencyInspectionTotal', label: '代办过检服务总价' },
    { key: 'twoCabinVipTotal', label: '两舱贵宾服务总价' },
    { key: 'businessVipTotal', label: '商务贵宾服务总价' },
    { key: 'memberCustomTotal', label: '会员定制服务总价' },
  ],
  materialFields: [
    { key: 'basicProduct', label: '基础产品' },
    { key: 'standardProduct', label: '标准产品' },
    { key: 'precisionProduct', label: '精密产品' },
    { key: 'intelligentProduct', label: '智能产品' },
  ],
  productionLineRows: [
    { key: 'changeProduct', label: '变更产品', labelColspan: 2 },
    { key: 'dismantleCost', label: '生产线拆除 1M / 条', labelColspan: 2 },
    { key: 'lineSale', label: '生产线出售', labelColspan: 2 },
    { key: 'newLineInstall', label: '新生产线安装', labelColspan: 2 },
    { key: 'constructionToFixed', label: '生产线残值', groupLabel: '建成生产线转固', groupRowspan: 2 },
    { key: 'newDepreciableAsset', label: '待折资产', useExistingGroupCell: true },
  ],
  researchFields: [
    { key: 'technologyResearch', label: '技术研发' },
    { key: 'managementSystem', label: '管理体系' },
  ],
  progressUpdateReminder: '产品进度更新',
} as const

const vipServiceReportLabels = {
  workInConstruction: '在建贵宾厅',
  factoryAsset: '贵宾区',
  lineResidual: '贵宾厅残值',
  workInProgress: '在途服务',
  finishedGoods: '服务验收',
  rawMaterials: '前期花费',
  row19Label: '3.服务进度更新',
  totalLiabilityEquity: '总负债权益',
  bestProductionHumanDirector: '最佳服务人力总监得分',
  taxRateEnterprise15: '0.15 支线机场贵宾',
} as const

const productionReportLabels = {
  workInConstruction: '在建生产线',
  factoryAsset: '厂房',
  lineResidual: '生产线残值',
  workInProgress: '在制品',
  finishedGoods: '成品',
  rawMaterials: '材料',
  row19Label: '税前利润',
  totalLiabilityEquity: '总负债和权益',
  bestProductionHumanDirector: '最佳生产人力总监得分',
  taxRateEnterprise15: '0.15 高新企业',
} as const

const vipServiceReportRequiredFieldLabels = {
  workInProgress: '在途服务',
  finishedGoods: '服务验收',
  rawMaterials: '前期花费',
  incomeTaxRate: '所得税税率',
  enterpriseCertificationScore: '企业认证得分',
  productionHumanScore: '最佳服务人力总监得分',
  closingSpeedScore: '关账速度得分',
} as const

const productionReportRequiredFieldLabels = {
  workInProgress: '在制品',
  finishedGoods: '成品',
  rawMaterials: '材料',
  incomeTaxRate: '所得税税率',
  enterpriseCertificationScore: '企业认证得分',
  productionHumanScore: '最佳生产人力总监得分',
  closingSpeedScore: '关账速度得分',
} as const

const vipServiceBaselineLabels = {
  workInConstruction: '在建贵宾厅',
  factoryAsset: '贵宾区',
  lineResidual: '贵宾厅残值',
  workInProgress: '在途服务',
  finishedGoods: '服务验收',
  rawMaterials: '前期花费',
} as const

const productionBaselineLabels = {
  workInConstruction: '在建生产线',
  factoryAsset: '厂房',
  lineResidual: '生产线残值',
  workInProgress: '在制品',
  finishedGoods: '成品',
  rawMaterials: '材料',
} as const

export type SandboxGameOperatingLabels = typeof vipServiceOperatingLabels | typeof productionOperatingLabels
export type SandboxGameReportLabels = typeof vipServiceReportLabels | typeof productionReportLabels
export type SandboxGameReportRequiredFieldLabels = typeof vipServiceReportRequiredFieldLabels | typeof productionReportRequiredFieldLabels
export type SandboxGameBaselineLabels = typeof vipServiceBaselineLabels | typeof productionBaselineLabels
export type DictionaryDisplayResolver = (itemCode: string, fallback: string) => string
type ProductionLineRowLabel = SandboxGameOperatingLabels['productionLineRows'][number] & {
  groupLabel?: string
}

export const serviceOperatingLabels = vipServiceOperatingLabels
export const serviceReportLabels = vipServiceReportLabels
export const serviceReportRequiredFieldLabels = vipServiceReportRequiredFieldLabels
export const serviceBaselineLabels = vipServiceBaselineLabels

export function resolveOperatingLabels(editionCode?: string | null): SandboxGameOperatingLabels {
  return isProductionEdition(editionCode) ? productionOperatingLabels : vipServiceOperatingLabels
}

export function resolveReportLabels(editionCode?: string | null): SandboxGameReportLabels {
  return isProductionEdition(editionCode) ? productionReportLabels : vipServiceReportLabels
}

export function resolveReportRequiredFieldLabels(editionCode?: string | null): SandboxGameReportRequiredFieldLabels {
  return isProductionEdition(editionCode) ? productionReportRequiredFieldLabels : vipServiceReportRequiredFieldLabels
}

export function resolveBaselineLabels(editionCode?: string | null): SandboxGameBaselineLabels {
  return isProductionEdition(editionCode) ? productionBaselineLabels : vipServiceBaselineLabels
}

export function applyDictionaryToOperatingLabels<T extends SandboxGameOperatingLabels>(labels: T, displayName: DictionaryDisplayResolver): T {
  return {
    ...labels,
    taxPolicyNote: displayName('operating.taxPolicyNote', labels.taxPolicyNote),
    materialGroupTitle: displayName('operating.materialGroupTitle', labels.materialGroupTitle),
    productionAdjustmentGroupTitle: displayName('operating.productionAdjustmentGroupTitle', labels.productionAdjustmentGroupTitle),
    humanResourceGroupTitle: displayName('operating.humanResourceGroupTitle', labels.humanResourceGroupTitle),
    humanResourceLabel: displayName('operating.humanResourceLabel', labels.humanResourceLabel),
    salaryGroupTitle: displayName('operating.salaryGroupTitle', labels.salaryGroupTitle),
    salaryLabel: displayName('operating.salaryLabel', labels.salaryLabel),
    researchAndManagementGroupTitle: displayName('operating.researchAndManagementGroupTitle', labels.researchAndManagementGroupTitle),
    newOrderReminder: displayName('operating.newOrderReminder', labels.newOrderReminder),
    receivableGroupTitle: displayName('operating.receivableGroupTitle', labels.receivableGroupTitle),
    receivableLabel: displayName('operating.receivableLabel', labels.receivableLabel),
    receivableReminder: displayName('operating.receivableReminder', labels.receivableReminder),
    deliveryGroupTitle: displayName('operating.deliveryGroupTitle', labels.deliveryGroupTitle),
    deliveryRevenueLabel: displayName('operating.deliveryRevenueLabel', labels.deliveryRevenueLabel),
    deliveryCostLabel: displayName('operating.deliveryCostLabel', labels.deliveryCostLabel),
    managementStaffGroupTitle: displayName('operating.managementStaffGroupTitle', labels.managementStaffGroupTitle),
    lineMaintenanceGroupTitle: displayName('operating.lineMaintenanceGroupTitle', labels.lineMaintenanceGroupTitle),
    lineMaintenanceNote: displayName('operating.lineMaintenanceNote', labels.lineMaintenanceNote),
    assetGroupTitle: displayName('operating.assetGroupTitle', labels.assetGroupTitle),
    assetValueNote: displayName('operating.assetValueNote', labels.assetValueNote),
    rentGroupTitle: displayName('operating.rentGroupTitle', labels.rentGroupTitle),
    rentValueNote: displayName('operating.rentValueNote', labels.rentValueNote),
    residualGroupTitle: displayName('operating.residualGroupTitle', labels.residualGroupTitle),
    depreciationGroupTitle: displayName('operating.depreciationGroupTitle', labels.depreciationGroupTitle),
    workInConstructionLabel: displayName('operating.workInConstructionLabel', labels.workInConstructionLabel),
    marketProductFields: labels.marketProductFields.map((item) => ({
      ...item,
      label: displayName(`operating.marketProductFields.${item.key}`, item.label),
    })),
    materialFields: labels.materialFields.map((item) => ({
      ...item,
      label: displayName(`operating.materialFields.${item.key}`, item.label),
    })),
    productionLineRows: labels.productionLineRows.map((rawItem) => {
      const item = rawItem as ProductionLineRowLabel
      return {
        ...item,
        label: displayName(`operating.productionLineRows.${item.key}`, item.label),
        groupLabel: item.groupLabel
          ? displayName(`operating.productionLineRows.${item.key}.groupLabel`, item.groupLabel)
          : item.groupLabel,
      }
    }),
    researchFields: labels.researchFields.map((item) => ({
      ...item,
      label: displayName(`operating.researchFields.${item.key}`, item.label),
    })),
    progressUpdateReminder: displayName('operating.progressUpdateReminder', labels.progressUpdateReminder),
  } as T
}

export function applyDictionaryToReportLabels<T extends SandboxGameReportLabels>(labels: T, displayName: DictionaryDisplayResolver): T {
  return {
    ...labels,
    workInConstruction: displayName('report.workInConstruction', labels.workInConstruction),
    factoryAsset: displayName('report.factoryAsset', labels.factoryAsset),
    lineResidual: displayName('report.lineResidual', labels.lineResidual),
    workInProgress: displayName('report.workInProgress', labels.workInProgress),
    finishedGoods: displayName('report.finishedGoods', labels.finishedGoods),
    rawMaterials: displayName('report.rawMaterials', labels.rawMaterials),
    row19Label: displayName('report.row19Label', labels.row19Label),
    totalLiabilityEquity: displayName('report.totalLiabilityEquity', labels.totalLiabilityEquity),
    bestProductionHumanDirector: displayName('report.bestProductionHumanDirector', labels.bestProductionHumanDirector),
    taxRateEnterprise15: displayName('report.taxRateEnterprise15', labels.taxRateEnterprise15),
  } as T
}

export function applyDictionaryToReportRequiredFieldLabels<T extends SandboxGameReportRequiredFieldLabels>(labels: T, displayName: DictionaryDisplayResolver): T {
  return {
    ...labels,
    workInProgress: displayName('reportRequired.workInProgress', labels.workInProgress),
    finishedGoods: displayName('reportRequired.finishedGoods', labels.finishedGoods),
    rawMaterials: displayName('reportRequired.rawMaterials', labels.rawMaterials),
    incomeTaxRate: displayName('reportRequired.incomeTaxRate', labels.incomeTaxRate),
    enterpriseCertificationScore: displayName('reportRequired.enterpriseCertificationScore', labels.enterpriseCertificationScore),
    productionHumanScore: displayName('reportRequired.productionHumanScore', labels.productionHumanScore),
    closingSpeedScore: displayName('reportRequired.closingSpeedScore', labels.closingSpeedScore),
  } as T
}

export function applyDictionaryToBaselineLabels<T extends SandboxGameBaselineLabels>(labels: T, displayName: DictionaryDisplayResolver): T {
  return {
    ...labels,
    workInConstruction: displayName('baseline.workInConstruction', labels.workInConstruction),
    factoryAsset: displayName('baseline.factoryAsset', labels.factoryAsset),
    lineResidual: displayName('baseline.lineResidual', labels.lineResidual),
    workInProgress: displayName('baseline.workInProgress', labels.workInProgress),
    finishedGoods: displayName('baseline.finishedGoods', labels.finishedGoods),
    rawMaterials: displayName('baseline.rawMaterials', labels.rawMaterials),
  } as T
}

function isProductionEdition(editionCode?: string | null) {
  return editionCode === 'PRODUCTION_V1'
}
