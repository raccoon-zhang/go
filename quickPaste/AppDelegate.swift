//
//  AppDelegate.swift
//  quickPaste
//
//  Created by zhangshaohua on 2024/9/25.
//

import Cocoa

@main
class AppDelegate: NSObject, NSApplicationDelegate {

    
    let statusItem = NSStatusBar.system.statusItem(withLength: NSStatusItem.variableLength)
    
    func applicationDidFinishLaunching(_ aNotification: Notification) {
        // Insert code here to initialize your application
        setupApp()
    }

    func applicationWillTerminate(_ aNotification: Notification) {
        // Insert code here to tear down your application
    }

    func applicationSupportsSecureRestorableState(_ app: NSApplication) -> Bool {
        return true
    }
}

extension AppDelegate {

    private func setupApp() {
        ClipBoard.shared.startListening()
        setupStatusBar()
    }

    private func setupStatusBar() {
        if let button = statusItem.button {
            button.image = NSImage(named: "StatusBarButtonImage")
            button.target = self
            button.action = #selector(statusMenuAction)
        }
    }
    
    @objc private func statusMenuAction() {
        let menu = NSMenu()
        let enableItem = NSMenuItem(title: "Enable", action: #selector(enableTool), keyEquivalent: "e")
        enableItem.state = ClipBoard.shared.enabled ? .on : .off
        let preFixItem = NSMenuItem()
        let itemView = PrefixView()
        itemView.textChangeBlock = { [weak self] prefix in
            guard self != nil else { return }
            CommonValues.shared.prefix = prefix
            ClipBoard.shared.onNewCopyFirst("prefix", { string in
                guard string.isUnderline() else {
                    return string
                }
                return prefix + string
            })
        }
        itemView.setupInitialValue(CommonValues.shared.prefix)
        preFixItem.view = itemView
        menu.addItem(enableItem)
        menu.addItem(preFixItem)
        menu.addItem(NSMenuItem(title: "Quit", action: #selector(NSApplication.terminate(_:)), keyEquivalent: "q"))
        menu.popUp(positioning: nil, at: NSMakePoint(0, statusItem.button?.frame.maxY ?? 0), in: statusItem.button)
    }
    
    @objc private func enableTool() {
        ClipBoard.shared.enabled = !ClipBoard.shared.enabled
    }
}

