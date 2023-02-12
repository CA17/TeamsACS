if (!window.settingsUi)
    window.settingsUi = {};


settingsUi.getConfigView = function (citem) {
    if (citem.name === "system") {
        return settingsUi.getSystemConfigView(citem);
    } else if (citem.name === "tr069") {
        return settingsUi.getTr069ConfigView(citem);
    }
    return {id: "settings_form_view"}
}


settingsUi.getSystemConfigView = function (citem) {
    let formid = webix.uid().toString();
    return {
        id: "settings_form_view",
        rows: [
            {
                padding: 7,
                cols: [
                    {
                        view: "label", label: " <i class='" + citem.icon + "'></i> " + citem.title,
                        css: "dash-title-b", width: 240, align: "left"
                    },
                    {},
                    wxui.getPrimaryButton(gtr("Save"), 150, false, function () {
                        let param = $$(formid).getValues();
                        param['ctype'] = 'system';
                        webix.ajax().post('/admin/settings/update', param).then(function (result) {
                            let resp = result.json();
                            webix.message({type: resp.msgtype, text: resp.msg, expire: 3000});
                        });
                    }),
                ],
            },
            {
                id: formid,
                view: "form",
                scroll: true,
                paddingX: 10,
                paddingY: 10,
                elementsConfig: {
                    labelWidth: 180,
                    labelPosition: "left",
                },
                url: "/admin/settings/system/query",
                elements: [
                    {
                        view: "radio", name: "SystemTheme", labelPosition: "top", label: tr("settings", "System Theme"),
                        options: ["light", "dark"]
                    },
                    {view: "text", name: "SystemTitle", labelPosition: "top", label: tr("settings", "Page title (browser title bar)")},
                    {view: "text", name: "SystemLoginRemark", labelPosition: "top", label: tr("settings", "Login screen prompt description")},
                    {view: "text", name: "SystemLoginSubtitle", labelPosition: "top", label: tr("settings", "Login form title")},
                    {}
                ],
            }
        ]
    }

}

settingsUi.getTr069ConfigView = function (citem) {
    let formid = webix.uid().toString();
    return {
        id: "settings_form_view",
        rows: [
            {
                padding: 7,
                cols: [
                    {
                        view: "label", label: " <i class='" + citem.icon + "'></i> " + citem.title,
                        css: "dash-title-b", width: 240, align: "left"
                    },
                    {},
                    wxui.getPrimaryButton(gtr("Save"), 150, false, function () {
                        let param = $$(formid).getValues();
                        param['ctype'] = 'tr069';
                        webix.ajax().post('/admin/settings/update', param).then(function (result) {
                            let resp = result.json();
                            webix.message({type: resp.msgtype, text: resp.msg, expire: 3000});
                        });
                    }),
                ],
            },
            {
                id: formid,
                view: "form",
                scroll: true,
                paddingX: 10,
                paddingY: 10,
                elementsConfig: {
                    labelWidth: 180,
                    labelPosition: "left",
                },
                url: "/admin/settings/tr069/query",
                elements: [
                    {
                        view: "radio", name: "CpeAutoRegister", labelPosition: "top", label: tr("settings", "Cpe auto register"),
                        options: ["enabled", "disabled"],
                        bottomLabel: tr("settings", "Automatic registration of new CPE devices")
                    },
                    {
                        view: "text", name: "TR069AccessAddress", labelPosition: "top", label: tr("settings", "TR069 access address"),
                        bottomLabel: tr("settings", "Teamsacs TR069 access address, HTTP | https://domain:port")
                    },
                    {
                        view: "text", name: "TR069AccessPassword", labelPosition: "top", label: tr("settings", "TR069 access password"),
                        bottomLabel: tr("settings", "Teamsacs TR069 access password, It is provided to CPE to access TeamsACS")
                    },
                    {
                        view: "text",
                        name: "CpeConnectionRequestPassword",
                        labelPosition: "top",
                        label: tr("settings", "CPE Connection authentication password"),
                        bottomLabel: tr("settings", "tr069 The authentication password used when the server connects to cpe")
                    },
                    {
                        view: "template", css: "form-desc", height: 200, borderless: true, src:"/admin/settings/tr069/quickset"
                    },
                ],
            }
        ]
    }
}
