import html
import zipfile
from datetime import datetime


OUTPUT = "QS2026_mapping_template.xlsx"

SUBJECTS = {
    "Business and Management Studies": "https://www.topuniversities.com/university-subject-rankings/business-management-studies",
    "Accounting and Finance": "https://www.topuniversities.com/university-subject-rankings/accounting-finance",
    "Statistics and Operational Research": "https://www.topuniversities.com/university-subject-rankings/statistics-operational-research",
    "Social Policy and Administration": "https://www.topuniversities.com/university-subject-rankings/social-policy-administration",
    "Law and Legal Studies": "https://www.topuniversities.com/university-subject-rankings/law-legal-studies",
    "Modern Languages": "https://www.topuniversities.com/university-subject-rankings/modern-languages",
    "Linguistics": "https://www.topuniversities.com/university-subject-rankings/linguistics",
    "Computer Science and Information Systems": "https://www.topuniversities.com/university-subject-rankings/computer-science-information-systems",
    "Data Science and Artificial Intelligence": "https://www.topuniversities.com/university-subject-rankings/data-science-artificial-intelligence",
    "Engineering - Electrical and Electronic": "https://www.topuniversities.com/university-subject-rankings/electrical-electronic-engineering",
    "Engineering - Mechanical": "https://www.topuniversities.com/university-subject-rankings/mechanical-aeronautical-manufacturing-engineering",
    "Engineering - Civil and Structural": "https://www.topuniversities.com/university-subject-rankings/civil-structural-engineering",
    "Mathematics": "https://www.topuniversities.com/university-subject-rankings/mathematics",
    "Art and Design": "https://www.topuniversities.com/university-subject-rankings/art-design",
    "Library and Information Management": "https://www.topuniversities.com/university-subject-rankings/library-information-management",
}

MAPPING = [
    ("工商管理", "Business and Management Studies", "首选", "QS 无工商管理细分，按商科管理大类查询"),
    ("行政管理", "Social Policy and Administration", "首选", "QS 无行政管理细分，公共事务/行政管理方向最接近"),
    ("汉语言文学", "Modern Languages", "首选", "QS 无汉语言文学细分，偏语言文化方向用此项"),
    ("汉语言文学", "Linguistics", "备选", "若偏语言学、语音、语义、应用语言学，可同时参考"),
    ("法学", "Law and Legal Studies", "首选", "QS 法学对应学科"),
    ("人力资源管理", "Business and Management Studies", "首选", "QS 无人力资源管理细分，按商科管理查询"),
    ("审计", "Accounting and Finance", "首选", "QS 无审计细分，按会计与金融查询"),
    ("金融", "Accounting and Finance", "首选", "QS 会计与金融对应学科"),
    ("统计", "Statistics and Operational Research", "首选", "QS 统计与运筹研究对应学科"),
    ("会计", "Accounting and Finance", "首选", "QS 会计与金融对应学科"),
    ("管理科学与工程", "Statistics and Operational Research", "首选", "偏管理科学、运筹优化、决策分析"),
    ("管理科学与工程", "Business and Management Studies", "备选", "若专业更偏管理学院/商科，可同时参考"),
    ("物流工程与管理", "Business and Management Studies", "首选", "QS 无物流/供应链独立学科，按商科管理查询"),
    ("物流工程与管理", "Statistics and Operational Research", "备选", "若偏物流优化、运筹、供应链建模，可同时参考"),
    ("供应链管理", "Business and Management Studies", "首选", "QS 无供应链独立学科，按商科管理查询"),
    ("供应链管理", "Statistics and Operational Research", "备选", "若偏供应链优化与运筹，可同时参考"),
    ("公共管理", "Social Policy and Administration", "首选", "公共管理/公共政策/行政方向最接近"),
    ("计算机类", "Computer Science and Information Systems", "首选", "计算机大类通用对应"),
    ("自动化类", "Engineering - Electrical and Electronic", "首选", "自动化/控制通常归入电气电子工程口径"),
    ("自动化类", "Data Science and Artificial Intelligence", "备选", "若方向偏智能控制/AI，可同时参考"),
    ("电子信息类", "Engineering - Electrical and Electronic", "首选", "电子信息、电子工程、通信相关通用对应"),
    ("电气类", "Engineering - Electrical and Electronic", "首选", "电气工程、电力系统、电气自动化通用对应"),
    ("机械类", "Engineering - Mechanical", "首选", "QS 官方该页对应 Mechanical, Aeronautical and Manufacturing Engineering"),
    ("土木类", "Engineering - Civil and Structural", "首选", "土木与结构工程对应学科"),
    ("通信工程", "Engineering - Electrical and Electronic", "首选", "通信工程通常归入电气电子工程"),
    ("通信工程", "Computer Science and Information Systems", "备选", "若课程更偏网络/信息系统，可同时参考"),
    ("软件工程", "Computer Science and Information Systems", "首选", "QS 无软件工程独立学科，按计算机与信息系统查询"),
    ("电子信息工程", "Engineering - Electrical and Electronic", "首选", "电子信息工程对应电气电子工程"),
    ("计算机科学与技术", "Computer Science and Information Systems", "首选", "计算机科学对应学科"),
    ("自动化", "Engineering - Electrical and Electronic", "首选", "自动化/控制工程方向最接近"),
    ("自动化", "Data Science and Artificial Intelligence", "备选", "若偏智能控制/机器人/AI，可同时参考"),
    ("控制工程", "Engineering - Electrical and Electronic", "首选", "控制工程通常归入电气电子工程"),
    ("电气工程及其自动化", "Engineering - Electrical and Electronic", "首选", "电气工程对应学科"),
    ("机械工程", "Engineering - Mechanical", "首选", "机械工程对应学科"),
    ("机电一体化", "Engineering - Mechanical", "首选", "QS 无机电一体化细分，机械工程最接近"),
    ("机电一体化", "Engineering - Electrical and Electronic", "备选", "若偏电气控制，可同时参考"),
    ("机械电子工程", "Engineering - Mechanical", "首选", "机械电子通常按机械工程查询"),
    ("机械电子工程", "Engineering - Electrical and Electronic", "备选", "若偏电子/控制，可同时参考"),
    ("电气自动化", "Engineering - Electrical and Electronic", "首选", "电气自动化对应电气电子工程"),
    ("网络工程", "Computer Science and Information Systems", "首选", "网络工程按计算机与信息系统查询"),
    ("网络工程", "Data Science and Artificial Intelligence", "备选", "若偏智能网络/安全算法，可参考"),
    ("物联网技术", "Computer Science and Information Systems", "首选", "QS 无物联网独立学科，按计算机与信息系统查询"),
    ("物联网技术", "Engineering - Electrical and Electronic", "备选", "若偏硬件/传感/通信，可同时参考"),
    ("物联网工程", "Computer Science and Information Systems", "首选", "QS 无物联网独立学科，按计算机与信息系统查询"),
    ("物联网工程", "Engineering - Electrical and Electronic", "备选", "若偏硬件/通信，可同时参考"),
    ("信息工程", "Computer Science and Information Systems", "首选", "信息工程按计算机/信息系统查询"),
    ("信息工程", "Engineering - Electrical and Electronic", "备选", "若偏电子通信，可同时参考"),
    ("人工智能", "Data Science and Artificial Intelligence", "首选", "QS 数据科学与人工智能对应学科"),
    ("人工智能", "Computer Science and Information Systems", "备选", "部分学校 AI 归在计算机学院，也可参考"),
    ("电子工程与信息技术", "Engineering - Electrical and Electronic", "首选", "电子工程对应电气电子工程"),
    ("电子工程与信息技术", "Computer Science and Information Systems", "备选", "若偏信息技术/软件，可同时参考"),
    ("计算机应用技术", "Computer Science and Information Systems", "首选", "计算机应用按计算机与信息系统查询"),
    ("计算机技术", "Computer Science and Information Systems", "首选", "计算机技术按计算机与信息系统查询"),
    ("计算机科学", "Computer Science and Information Systems", "首选", "计算机科学对应学科"),
    ("信息安全", "Computer Science and Information Systems", "首选", "QS 无信息安全独立学科，按计算机与信息系统查询"),
    ("信息安全", "Data Science and Artificial Intelligence", "备选", "若偏 AI 安全/数据安全，可参考"),
    ("应用数学", "Mathematics", "首选", "应用数学按数学查询"),
    ("电子科学与技术", "Engineering - Electrical and Electronic", "首选", "电子科学与技术对应电气电子工程"),
    ("网络与信息安全", "Computer Science and Information Systems", "首选", "网络安全/信息安全按计算机与信息系统查询"),
    ("软件工程与智能系统", "Computer Science and Information Systems", "首选", "软件工程按计算机与信息系统查询"),
    ("软件工程与智能系统", "Data Science and Artificial Intelligence", "备选", "智能系统方向可参考 AI 学科"),
    ("计算机软件", "Computer Science and Information Systems", "首选", "软件方向按计算机与信息系统查询"),
    ("信息科技技术", "Computer Science and Information Systems", "首选", "信息技术按计算机与信息系统查询"),
    ("信息管理", "Library and Information Management", "首选", "若偏信息资源/信息管理，可查该学科"),
    ("信息管理", "Computer Science and Information Systems", "备选", "若偏信息系统/IT 管理，可同时参考"),
    ("人机交互", "Computer Science and Information Systems", "首选", "QS 无人机交互独立学科，技术口径按计算机查询"),
    ("人机交互", "Art and Design", "备选", "若偏交互设计/用户体验，可同时参考"),
    ("交互设计", "Art and Design", "首选", "交互设计按艺术与设计查询"),
    ("交互设计", "Computer Science and Information Systems", "备选", "若偏 HCI 技术实现，可同时参考"),
    ("工业设计", "Art and Design", "首选", "工业设计按艺术与设计查询"),
    ("数字媒体技术", "Art and Design", "首选", "若偏数字媒体/设计作品集，可查艺术与设计"),
    ("数字媒体技术", "Computer Science and Information Systems", "备选", "若偏技术开发/图形学，可同时参考"),
    ("电力系统及其自动化", "Engineering - Electrical and Electronic", "首选", "电力系统归入电气电子工程"),
    ("电气自动化技术", "Engineering - Electrical and Electronic", "首选", "电气自动化技术对应电气电子工程"),
    ("机电一体化技术", "Engineering - Mechanical", "首选", "机电一体化技术按机械工程查询"),
    ("机电一体化技术", "Engineering - Electrical and Electronic", "备选", "若偏电气控制，可同时参考"),
    ("计算机/软件工程", "Computer Science and Information Systems", "首选", "计算机/软件工程按计算机与信息系统查询"),
]


def col_name(n):
    name = ""
    while n:
        n, r = divmod(n - 1, 26)
        name = chr(65 + r) + name
    return name


def sheet_xml(rows):
    out = [
        '<?xml version="1.0" encoding="UTF-8" standalone="yes"?>',
        '<worksheet xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main">',
        "<sheetData>",
    ]
    for r_idx, row in enumerate(rows, 1):
        out.append(f'<row r="{r_idx}">')
        for c_idx, val in enumerate(row, 1):
            ref = f"{col_name(c_idx)}{r_idx}"
            text = html.escape(str(val) if val is not None else "")
            out.append(f'<c r="{ref}" t="inlineStr"><is><t>{text}</t></is></c>')
        out.append("</row>")
    out.append("</sheetData></worksheet>")
    return "".join(out)


def build_workbook():
    unique_subjects = []
    seen = set()
    for _, subject, _, _ in MAPPING:
        if subject not in seen:
            seen.add(subject)
            unique_subjects.append(subject)

    sheets = [
        (
            "专业_QS学科映射",
            [["原始专业", "QS学科", "映射优先级", "QS官网查询页", "备注"]]
            + [[major, subject, priority, SUBJECTS[subject], note] for major, subject, priority, note in MAPPING],
        ),
        (
            "QS学科去重清单",
            [["QS学科", "QS官网查询页", "抓取范围", "建议抓取字段", "备注"]]
            + [
                [
                    subject,
                    SUBJECTS[subject],
                    "前50",
                    "Rank / Institution 或 University",
                    "打开页面后选择 2026；若有 Download Excel Table 可优先下载",
                ]
                for subject in unique_subjects
            ],
        ),
        (
            "抓取结果模板",
            [["原始专业", "QS学科", "QS排名", "院校英文名", "院校中文名"]]
            + [[major, subject, "", "", ""] for major, subject, priority, _ in MAPPING if priority == "首选"],
        ),
        (
            "插件操作简表",
            [
                ["步骤", "操作说明"],
                ["1", "打开 QS学科去重清单 中的官网查询页，确认标题为 QS World University Rankings by Subject 2026。"],
                ["2", "若 QS 页面有 Download Excel Table，优先用官网下载；否则再用 Table Capture 或 Instant Data Scraper。"],
                ["3", "只抓 Rank 和 Institution/University 两列；复制到 抓取结果模板 的 QS排名、院校英文名。"],
                ["4", "同一 QS 学科只需抓一次，多个原始专业共用同一份前 50 名单。"],
                ["5", "院校中文名可以先空着，英文名抓完后再批量补中文译名。"],
            ],
        ),
    ]

    content_types = [
        '<?xml version="1.0" encoding="UTF-8" standalone="yes"?>',
        '<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">',
        '<Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>',
        '<Default Extension="xml" ContentType="application/xml"/>',
        '<Override PartName="/xl/workbook.xml" ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.sheet.main+xml"/>',
        '<Override PartName="/docProps/core.xml" ContentType="application/vnd.openxmlformats-package.core-properties+xml"/>',
        '<Override PartName="/docProps/app.xml" ContentType="application/vnd.openxmlformats-officedocument.extended-properties+xml"/>',
    ]
    for i in range(1, len(sheets) + 1):
        content_types.append(
            f'<Override PartName="/xl/worksheets/sheet{i}.xml" ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.worksheet+xml"/>'
        )
    content_types.append("</Types>")

    workbook_sheets = "".join(
        f'<sheet name="{html.escape(name)}" sheetId="{i}" r:id="rId{i}"/>'
        for i, (name, _) in enumerate(sheets, 1)
    )
    workbook = (
        '<?xml version="1.0" encoding="UTF-8" standalone="yes"?>'
        '<workbook xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main" '
        f'xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships"><sheets>{workbook_sheets}</sheets></workbook>'
    )

    workbook_rels = [
        '<?xml version="1.0" encoding="UTF-8" standalone="yes"?>',
        '<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">',
    ]
    for i in range(1, len(sheets) + 1):
        workbook_rels.append(
            f'<Relationship Id="rId{i}" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/worksheet" Target="worksheets/sheet{i}.xml"/>'
        )
    workbook_rels.append("</Relationships>")

    now = datetime.utcnow().strftime("%Y-%m-%dT%H:%M:%SZ")
    core = (
        '<?xml version="1.0" encoding="UTF-8" standalone="yes"?>'
        '<cp:coreProperties xmlns:cp="http://schemas.openxmlformats.org/package/2006/metadata/core-properties" '
        'xmlns:dc="http://purl.org/dc/elements/1.1/" xmlns:dcterms="http://purl.org/dc/terms/" '
        'xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">'
        f"<dc:creator>Codex</dc:creator><cp:lastModifiedBy>Codex</cp:lastModifiedBy>"
        f'<dcterms:created xsi:type="dcterms:W3CDTF">{now}</dcterms:created>'
        f'<dcterms:modified xsi:type="dcterms:W3CDTF">{now}</dcterms:modified>'
        "</cp:coreProperties>"
    )
    app = (
        '<?xml version="1.0" encoding="UTF-8" standalone="yes"?>'
        '<Properties xmlns="http://schemas.openxmlformats.org/officeDocument/2006/extended-properties" '
        'xmlns:vt="http://schemas.openxmlformats.org/officeDocument/2006/docPropsVTypes">'
        "<Application>Codex</Application></Properties>"
    )
    root_rels = (
        '<?xml version="1.0" encoding="UTF-8" standalone="yes"?>'
        '<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">'
        '<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="xl/workbook.xml"/>'
        '<Relationship Id="rId2" Type="http://schemas.openxmlformats.org/package/2006/relationships/metadata/core-properties" Target="docProps/core.xml"/>'
        '<Relationship Id="rId3" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/extended-properties" Target="docProps/app.xml"/>'
        "</Relationships>"
    )

    with zipfile.ZipFile(OUTPUT, "w", zipfile.ZIP_DEFLATED) as z:
        z.writestr("[Content_Types].xml", "".join(content_types))
        z.writestr("_rels/.rels", root_rels)
        z.writestr("xl/workbook.xml", workbook)
        z.writestr("xl/_rels/workbook.xml.rels", "".join(workbook_rels))
        z.writestr("docProps/core.xml", core)
        z.writestr("docProps/app.xml", app)
        for i, (_, rows) in enumerate(sheets, 1):
            z.writestr(f"xl/worksheets/sheet{i}.xml", sheet_xml(rows))


if __name__ == "__main__":
    build_workbook()
    print(OUTPUT)
